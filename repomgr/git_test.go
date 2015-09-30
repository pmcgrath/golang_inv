package main

import (
	"os"
	"testing"
)

func TestFilterGitReposOnly(t *testing.T) {
	currentDirWhichIsAGitRepo, _ := os.Getwd()
	tempDirWhichIsNotAGitRepo := os.TempDir()

	for _, testCase := range []struct {
		Path           string
		ExpectedResult bool
	}{
		{currentDirWhichIsAGitRepo, true},
		{tempDirWhichIsNotAGitRepo, false},
	} {
		gitRepos := filterGitReposOnly([]string{testCase.Path})
		t.Logf("***FOR [%s] gitRepos = %s\n", testCase.Path, gitRepos)
		actualResult := len(gitRepos) == 1 && gitRepos[0] == testCase.Path
		if actualResult != testCase.ExpectedResult {
			t.Errorf("Unexpected result for [%s], expected %t but got %t", testCase.Path, testCase.ExpectedResult, actualResult)
		}
	}
}
