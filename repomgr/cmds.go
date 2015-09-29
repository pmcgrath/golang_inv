package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func execCmd(cmdName string, args ...string) ([]string, error) {
	cmd := exec.Command(cmdName, args...)
	cmdOutputBytes, err := cmd.CombinedOutput()
	cmdOutputString := string(cmdOutputBytes)
	if strings.TrimSpace(cmdOutputString) == "" {
		// No output (No error) or a raw error
		return nil, err
	}

	if err != nil {
		// Use output as error
		return nil, fmt.Errorf(cmdOutputString)
	}

	outputLines := strings.Split(cmdOutputString, "\n")

	return outputLines, err
}
