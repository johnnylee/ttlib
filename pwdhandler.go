package ttlib

import (
	"fmt"
	"github.com/johnnylee/util"
	"path/filepath"
	"sync"
	"time"
)

// PwdHandler provides an interface for requesting user's passwords.
type PwdHandler struct {
	mutex   sync.RWMutex
	hash    map[string][]byte
	baseDir string
}

func NewPwdHandler(baseDir string) *PwdHandler {
	ph := new(PwdHandler)
	ph.baseDir = baseDir
	ph.hash = make(map[string][]byte)
	go ph.fileWatcher()
	return ph
}

func (ph *PwdHandler) GetPwdHash(user string) ([]byte, error) {
	ph.mutex.RLock()
	defer ph.mutex.RUnlock()
	hash, ok := ph.hash[user]
	if !ok {
		return nil, fmt.Errorf("Unknown user: %v", user)
	}
	return hash, nil
}

// fileWatcher: a state machine that watches for new or removed files. Check
// using gotofsm-verify.
func (ph *PwdHandler) fileWatcher() {
addUsers:
	ph.addUsers()
	goto deleteUsers

deleteUsers:
	ph.deleteUsers()
	goto sleep

sleep:
	time.Sleep(60 * time.Second)
	goto addUsers
}

// addUsers: Used by fileWatcher.
func (ph *PwdHandler) addUsers() {
	pattern := filepath.Join(ClientPwdDir(ph.baseDir), "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Err(err, "Error reading pwd file list")
		return
	}

	for _, path := range matches {
		user := filepath.Base(path)
		user = user[:len(user)-5]

		// Skip known paths.
		if _, ok := ph.hash[user]; ok {
			continue
		}

		// Open the password file.
		pwdFile, err := LoadPwdFile(ph.baseDir, user)
		if err != nil {
			log.Err(err, "Error loading password file for user: %v", user)
			continue
		}

		log.Msg("Adding user: %v", user)

		ph.mutex.Lock()
		ph.hash[user] = pwdFile.PwdHash
		ph.mutex.Unlock()
	}
}

// deleteUsers: Used by fileWatcher.
func (ph *PwdHandler) deleteUsers() {
	for user := range ph.hash {
		if !util.FileExists(ClientPwdFile(ph.baseDir, user)) {
			log.Msg("Removing user: %v", user)
			ph.mutex.Lock()
			delete(ph.hash, user)
			ph.mutex.Unlock()
		}
	}

}
