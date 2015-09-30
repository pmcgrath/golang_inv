package main

import (
	"os"
	"path"
	"testing"
)

func TestExecCmd(t *testing.T) {
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
		t.Logf("---> %s: %t\n", testCase.Path, errEncountered)
		if errEncountered != testCase.ExpectedError {
			t.Errorf("Unexpected result for [%s], expected to encounter an error = %t but got an error = %t", testCase.Path, testCase.ExpectedError, errEncountered)
		}
	}
}
