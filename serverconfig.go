package ttlib

import (
	"github.com/johnnylee/util"
)

// ServerConfig: A configuration file for a server.
type ServerConfig struct {
	ListenAddr string // The address to pass to Listen.
	PublicAddr string // The address to use to connect to the server.
}

// LoadServerConfig: Load the server's configuration from the given baseDir.
func LoadServerConfig(baseDir string) (*ServerConfig, error) {
	sc := new(ServerConfig)
	return sc, util.JsonUnmarshal(ServerConfigPath(baseDir), sc)
}

// Save: Save the server's configuration to the given baseDir.
func (sc ServerConfig) Save(baseDir string) error {
	return util.JsonMarshal(ServerConfigPath(baseDir), &sc)
}
