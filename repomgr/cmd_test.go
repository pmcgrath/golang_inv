package main

import (
	"os"
	"path"
	"testing"
)

func TestExecCmdForKnownCommand(t *testing.T) {
	output, err := execCmd("hostname")
	if err != nil {
		t.Errorf("Unexpected result for running a known command, error: %#v", err)
	}
	if len(output) < 1 || len(output[0]) < 1 {
		t.Error("Expected a result for running a known command")
	}
}

func TestExecCmdForUnknownCommand(t *testing.T) {
	_, err := execCmd("ted_unknown")
	if err == nil {
		t.Error("Unexpected result for running unkown command")
	}
}

func TestExecCmdForLs(t *testing.T) {
	currentDir, _ := os.Getwd()
	tempDir := os.TempDir()
	nonExistentDir := path.Join(tempDir, "does_not_exist")

	for _, testCase := range []struct {
		Path          string
		ExpectedError bool
	}{
		{currentDir, false},
		{tempDir, false},
		{nonExistentDir, true},
	} {
		_, err := execCmd("ls", "-al", testCase.Path)
		errEncountered := err != nil
		if errEncountered != testCase.ExpectedError {
			t.Errorf("Unexpected result for [%s], expected to encounter an error = %t but got an error = %t", testCase.Path, testCase.ExpectedError, errEncountered)
		}
	}
}
