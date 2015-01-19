package ttlib

import (
	"github.com/johnnylee/util"
)

// PwdFile: A password file (on the server) for a client.
type PwdFile struct {
	PwdHash []byte
}

// LoadPwdFile: Load a password file for the given user.
func LoadPwdFile(baseDir, user string) (*PwdFile, error) {
	pf := new(PwdFile)
	return pf, util.JsonUnmarshal(ClientPwdFile(baseDir, user), pf)
}

// Save: Save the password file for the given user.
func (pf PwdFile) Save(baseDir, user string) error {
	return util.JsonMarshal(ClientPwdFile(baseDir, user), &pf)
}
