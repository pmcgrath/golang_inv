package main

import (
	"os/user"
	"strings"
)

func getCurrentUserName() (name string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	name = strings.ToLower(user.Username)
	domainSeperatorIndex := strings.Index(name, "\\")
	if domainSeperatorIndex > -1 {
		name = name[domainSeperatorIndex+1:]
	}

	return
}
