package ttlib

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

// Dial is a wrapper for tls.Dial. It uses a configuration file containing
// the remote address, username, password, and certificate.
func Dial(configPath string) (*tls.Conn, error) {
	// Load the configuration file.
	cc, err := LoadClientConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool for the client.
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cc.CaCert) {
		return nil, fmt.Errorf("Error appending certificate to pool.")
	}

	// Create TLS configuration.
	config := tls.Config{RootCAs: certPool}

	// Dial out.
	conn, err := tls.Dial("tcp", cc.Host, &config)
	if err != nil {
		return nil, err
	}

	// Send name.
	if err = writeBytes(conn, []byte(cc.User)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Send password.
	if err = writeBytes(conn, cc.Pwd); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Return the connection.
	return conn, nil
}

// Client: A wrapper for rpc.Client that performs automatic reconnects.
type RpcClient struct {
	configPath string
	mutex      sync.Mutex  // Protects access to client.
	client     *rpc.Client // The rpc client pointer.
}

func RpcDial(configPath string) *RpcClient {
	rc := new(RpcClient)
	rc.configPath = configPath
	rc.client = nil
	return rc
}

// Get the current client.
func (rc *RpcClient) getClient() *rpc.Client {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	return rc.client
}

// Reconnect if the current client is different from the given client.
func (rc *RpcClient) reconnect(client *rpc.Client) (*rpc.Client, error) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// If the current client isn't nil and is new to the caller, return it.
	if rc.client != nil && client != rc.client {
		return rc.client, nil
	}

	// Create a new connection to the server.
	conn, err := Dial(rc.configPath)
	if err != nil {
		// Sleep for a small amount of time here to slow down reconnect
		// attempts to the server.
		log.Err(err, "Failed to connect to server (sleeping 4 seconds)")
		time.Sleep(4 * time.Second)
		rc.client = nil
		return nil, err
	}

	log.Msg("Connected to server.")

	// Construct a new client using the encrypted connection.
	rc.client = rpc.NewClient(conn)
	return rc.client, nil
}

func (rc *RpcClient) Call(
	serviceMethod string, args interface{}, reply interface{}) error {

	var err error

	// Get client.
	client := rc.getClient()

	// Attempt reconnect if client is nil.
	if client == nil {
		if client, err = rc.reconnect(client); err != nil {
			return err
		}
	}

	// First attempt at call.
	err = client.Call(serviceMethod, args, reply)

	// If the call succeeded, then return the error code. The only code not
	// returned is rpc.ErrShutdown, which means that the connection to the
	// server has been lost.
	if err != rpc.ErrShutdown {
		return err
	}

	// Attempt reconnect.
	if client, err = rc.reconnect(client); err != nil {
		return err
	}

	// Make a final attempt at the function call.
	return client.Call(serviceMethod, args, reply)
}

func (rc *RpcClient) Close() error {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	err := rc.client.Close()
	rc.client = nil
	return err
}
