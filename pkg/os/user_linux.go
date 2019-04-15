package os

import (
	log "github.com/code-ready/crc/pkg/crc/logging"
	"os/user"
	"strconv"
)

func CheckUserPrivilages() bool {
	cuser, err := user.Current()
	if err != nil {
		log.Error(err)
		return false
	}

	id, err := strconv.Atoi(cuser.Uid)
	if err != nil {
		log.Error(err)
		return false
	}
	if id == 0 {
		return true
	}

	return false
}
