package main

import (
	"log"
	"os/user"
	"strings"
)

func getCurrentUserName() (name string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	log.Printf("OS username: %s\n", user.Username)
	name = strings.ToLower(user.Username)
	domainSeperatorIndex := strings.Index(name, "\\")
	if domainSeperatorIndex > -1 {
		name = name[domainSeperatorIndex+1:]
	}

	return
}

func getCurrentUserHomedir() (dir string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	dir = user.HomeDir

	return
}
