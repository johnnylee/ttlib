package ttlib

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/tls"
	"github.com/johnnylee/util"
	"net"
	"net/rpc"
)

type acceptResult struct {
	conn net.Conn
	err  error
}

/* We construct TtListener so that it accept connections in a go routine, and
 * sends them back to the calling function after they've been authenticated.
 */

type TtListener struct {
	ln         net.Listener
	acceptChan chan acceptResult
	pwdHandler *PwdHandler
}

// Accept wrapper.
func (ln *TtListener) Accept() (net.Conn, error) {
	go func() {
		var name []byte
		var hash []byte
		var pwd []byte

		// Accept the actual connection.
		conn, err := ln.ln.Accept()
		if err != nil {
			goto finished
		}

		// Get the client's name.
		name, err = readBytes(conn, 256)
		if err != nil {
			goto finished
		}

		// Get the client's password hash.
		hash, err = ln.pwdHandler.GetPwdHash(string(name))
		if err != nil {
			goto finished
		}

		// Read the client's password.
		pwd, err = readBytes(conn, 256)
		if err != nil {
			goto finished
		}

		// Check that the password is correct.
		if err = bcrypt.CompareHashAndPassword(hash, pwd); err != nil {
			goto finished
		}

	finished:
		if err != nil {
			_ = conn.Close()
		} else {
			log.Msg("Accepting connection for user: %v", string(name))
		}
		ln.acceptChan <- acceptResult{conn, err}
	}()

	// Return connection from channel.
	result := <-ln.acceptChan
	return result.conn, result.err
}

func (ln *TtListener) Close() error {
	return ln.ln.Close()
}

func (ln *TtListener) Addr() net.Addr {
	return ln.ln.Addr()
}

// Listen is a wrapper for tls.Listen.
func Listen(baseDir string) (*TtListener, error) {

	var err error
	if baseDir, err = util.ExpandPath(baseDir); err != nil {
		return nil, err
	}

	// Read in the configuration file.
	serverConfig, err := LoadServerConfig(baseDir)
	if err != nil {
		return nil, err
	}

	// Load key and certificate.
	cert, err := tls.LoadX509KeyPair(
		ServerCertPath(baseDir), ServerKeyPath(baseDir))
	if err != nil {
		return nil, err
	}

	// Create TLS configuration.
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	// Listen for connections
	ln := new(TtListener)
	ln.ln, err = tls.Listen("tcp", serverConfig.ListenAddr, &config)
	if err != nil {
		return nil, err
	}

	// Create the accept-result channel.
	ln.acceptChan = make(chan acceptResult)

	// Create the password handler.
	ln.pwdHandler = NewPwdHandler(baseDir)

	return ln, nil
}

func RpcServeForever(baseDir string, obj interface{}) error {
	// Register the server to handle function calls.
	if err := rpc.Register(obj); err != nil {
		log.Err(err, "Failed to register object")
		return err
	}

	// Start the network listener.
	ln, err := Listen(baseDir)
	if err != nil {
		log.Err(err, "Failed to listen")
		return err
	}

	// Accept connections forever.
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Err(err, "Failed to accept connection")
			continue
		}
		go rpc.ServeConn(conn)
	}
}
