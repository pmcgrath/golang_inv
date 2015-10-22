package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExecAllDirectGitCommandsExcludingClone(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := make([]string, 1)
	for index := 0; index < len(repoPaths); index++ {
		repoPaths[index] = fmt.Sprintf("/tmp/repos/repo%03d", index+1)
	}
	remoteName := "ARemoteSomewhere"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	for _, cmd := range []command{
		branchCmd,
		fetchCmd,
		pullCmd,
		remoteCmd,
		statusCmd} {

		var cmdSuffix string
		var results gitCmdResults
		switch cmd {
		case branchCmd:
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 branch -av
			// So suffix is "-av"
			cmdSuffix = "-av"
			results = execGitBranch(repoPaths)
		case fetchCmd:
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 fetch origin
			// So suffix is "origin" or whatever remote name we are running against
			cmdSuffix = remoteName
			results = execGitFetch(repoPaths, remoteName)
		case pullCmd:
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 pull origin
			// So suffix is "origin" or whatever remote name we are running against
			cmdSuffix = remoteName
			results = execGitPull(repoPaths, remoteName)
		case remoteCmd:
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 remote --verbose
			// So suffix is "--verbose"
			cmdSuffix = "--verbose"
			results = execGitRemote(repoPaths)
		case statusCmd:
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 status --porcelain
			// So suffix is "--porcelain"
			cmdSuffix = "--porcelain"
			results = execGitStatus(repoPaths)
		}

		if len(results) != len(repoPaths) {
			t.Errorf("Unexpected result count : %d", len(results))
		}
		for index, repoPath := range repoPaths {
			repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

			if repoResult.RepoPath != repoPath {
				t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
			}
			if repoResult.Command != string(cmd) {
				t.Errorf("Unexpected result Command, expected %s, but got %s", cmd, repoResult.Command)
			}
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 branch -av
			expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s %s %s", repoPath, repoPath, cmd, cmdSuffix)
			if repoResult.Output[0] != expectedOutput {
				t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
			}
		}
	}
}

func TestExecGitBranchWithAFailure(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := make([]string, 100)
	for index := 0; index < len(repoPaths); index++ {
		repoPaths[index] = fmt.Sprintf("/tmp/repos/repo%03d", index+1)
	}
	erroredRepoPath := "/tmp/repos/repo011"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	execCmdErr := fmt.Errorf("BANG")
	cmdExecutionFn = func(cmdName string, args ...string) ([]string, error) {
		if args[3] == erroredRepoPath {
			return nil, execCmdErr
		}
		return execCmdSpy(cmdName, args...)
	}

	results := execGitBranch(repoPaths)

	if len(results) != len(repoPaths) {
		t.Errorf("Unexpected result count : %d", len(results))
	}
	for index, repoPath := range repoPaths {
		repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

		if repoResult.RepoPath != repoPath {
			t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
		}
		if repoResult.Command != branchCmd {
			t.Errorf("Unexpected result Command, expected %s, but got %s", branchCmd, repoResult.Command)
		}
		if repoResult.RepoPath == erroredRepoPath {
			if repoResult.Error != execCmdErr {
				t.Errorf("Unexpected result Error, expected %s, but got %s", execCmdErr, repoResult.Error)
			}
		} else {
			// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 branch -av
			expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s %s -av", repoPath, repoPath, branchCmd)
			if repoResult.Output[0] != expectedOutput {
				t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
			}
		}
	}
}

func TestExecGitClone(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoURLs := []string{"git@bitbucket.org:pmcgrath/bbdotfiles.git", "git@github.com:pmcgrath/dotfiles.git", "git@github.com:pmcgrath/other.git"}
	remoteName := "some_other_remote"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitClone("/tmp", repoURLs, remoteName)
	if len(results) != len(repoURLs) {
		t.Errorf("Unexpected result count : %d", len(results))
	}
	for index, repoURL := range repoURLs {
		repoResult := results[index] // Since we ensured the input is sorted, we can rely on the order

		if repoResult.RepoPath != repoURL {
			t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoURL, repoResult.RepoPath)
		}
		if repoResult.Command != cloneCmd {
			t.Errorf("Unexpected result Command, expected %s, but got %s", cloneCmd, repoResult.Command)
		}
		// Sample command is : git clone --origin github git@github.com:pmcgrath/dotfiles.git
		expectedOutput := fmt.Sprintf("cmdName: git, args: clone --origin %s %s", remoteName, repoURL)
		if repoResult.Output[0] != expectedOutput {
			t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
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
