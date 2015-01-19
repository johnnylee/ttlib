package ttlib

import (
	"github.com/johnnylee/util"
)

// ClientConfig: A configuration file for a client.
type ClientConfig struct {
	Host   string // The host address: <address>:<port>.
	User   string // The username for the client.
	Pwd    []byte // The user's password.
	CaCert []byte // The CA certificate.
}

// Load: Load the client's configuration from the given path.
func LoadClientConfig(path string) (*ClientConfig, error) {
	cc := new(ClientConfig)
	return cc, util.JsonUnmarshal(path, cc)
}

// Save: Save the client's configuration to the given path.
func (cc ClientConfig) Save(path string) error {
	return util.JsonMarshal(path, &cc)
}
