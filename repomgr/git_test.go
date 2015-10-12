package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExecGitBranch(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := []string{"/tmp/repos/repo1", "/tmp/repos/repo2"}

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitBranch(repoPaths)

	if len(results) != len(repoPaths) {
		t.Errorf("Unexpected result count : %d", len(results))
	}
	for index, repoPath := range repoPaths {
		repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

		if repoResult.RepoPath != repoPath {
			t.Logf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
		}
		if repoResult.Command != "branch" {
			t.Logf("Unexpected result Command, expected %s, but got %s", "branch", repoResult.Command)
		}
		// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 branch -av
		expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s branch -av", repoPath, repoPath)
		if repoResult.Output[0] != expectedOutput {
			t.Logf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
		}
	}
}

func TestExecGitClone(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoUrls := []string{"git@github.com:pmcgrath/dotfiles.git"}
	remoteName := "mygithub"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitClone("/tmp", repoUrls, remoteName)
	if len(results) != len(repoUrls) {
		t.Errorf("Unexpected result count : %d", len(results))
	}
	for index, repoUrl := range repoUrls {
		repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

		if repoResult.RepoPath != repoUrl {
			t.Logf("Unexpected result RepoPath, expected %s, but got %s", repoUrl, repoResult.RepoPath)
		}
		if repoResult.Command != "clone" {
			t.Logf("Unexpected result Command, expected %s, but got %s", "clone", repoResult.Command)
		}
		// Sample command is : git clone --origin github git@github.com:pmcgrath/dotfiles.git
		expectedOutput := fmt.Sprintf("cmdName: git, args: clone --origin %s %s", remoteName, repoUrl)
		if repoResult.Output[0] != expectedOutput {
			t.Logf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
		}
	}
}

func TestExecGitFetch(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := []string{"/tmp/repos/repo1", "/tmp/repos/repo2"}
	remoteName := "origin"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitFetch(repoPaths, remoteName)

	if len(results) != len(repoPaths) {
		t.Errorf("Unexpected result count : %d", len(results))
	}
	for index, repoPath := range repoPaths {
		repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

		if repoResult.RepoPath != repoPath {
			t.Logf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
		}
		if repoResult.Command != "fetch" {
			t.Logf("Unexpected result Command, expected %s, but got %s", "branch", repoResult.Command)

		}
		// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 fetch origin
		expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s fetch %s", repoPath, repoPath, remoteName)
		if repoResult.Output[0] != expectedOutput {
			t.Logf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
		}
	}
}

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
		actualResult := len(gitRepos) == 1 && gitRepos[0] == testCase.Path
		if actualResult != testCase.ExpectedResult {
			t.Errorf("Unexpected result for [%s], expected %t but got %t", testCase.Path, testCase.ExpectedResult, actualResult)
		}
	}
}

func execCmdSpy(cmdName string, args ...string) ([]string, error) {
	return []string{fmt.Sprintf("cmdName: %s, args: %s", cmdName, strings.Join(args, " "))}, nil
}
