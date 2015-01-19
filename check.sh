#!/bin/bash

go vet &&
gotofsm-verify dot/filewatcher.dot pwdhandler.go fileWatcher &&
errcheck github.com/johnnylee/ttlib
