// +build linux

package main

import (
	"os/exec"
)

func IsProcessRunning(programName string) bool {
	// This is linux specific - using build tags - http://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool
	cmd := exec.Command("/bin/pidof", programName)
	buf, _ := cmd.Output()
	return len(buf) > 0
}
