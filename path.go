package ttlib

import (
	"path/filepath"
)

func ServerConfigPath(baseDir string) string {
	return filepath.Join(baseDir, "config.json")
}

func ServerKeyPath(baseDir string) string {
	return filepath.Join(baseDir, "server.key")
}

func ServerCertPath(baseDir string) string {
	return filepath.Join(baseDir, "server.crt")
}

func ClientPwdDir(baseDir string) string {
	return filepath.Join(baseDir, "clients")
}

func ClientPwdFile(baseDir, user string) string {
	return filepath.Join(baseDir, "clients", user+".json")
}

func ClientConfigDir(baseDir string) string {
	return filepath.Join(baseDir, "client-config")
}

func ClientConfigFile(baseDir, user string) string {
	return filepath.Join(baseDir, "client-config", user+".json")
}
