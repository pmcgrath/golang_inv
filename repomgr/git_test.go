package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExecGitBranch(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := make([]string, 100)
	for index := 0; index < len(repoPaths); index++ {
		repoPaths[index] = fmt.Sprintf("/tmp/repos/repo%03d", index+1)
	}

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
			t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
		}
		if repoResult.Command != branchCmd {
			t.Errorf("Unexpected result Command, expected %s, but got %s", branchCmd, repoResult.Command)
		}
		// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 branch -av
		expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s branch -av", repoPath, repoPath)
		if repoResult.Output[0] != expectedOutput {
			t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
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
	execCmdErr := fmt.Errorf("BANG!")
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
			expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s branch -av", repoPath, repoPath)
			if repoResult.Output[0] != expectedOutput {
				t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
			}
		}
	}
}

func TestExecGitClone(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoUrls := []string{"git@bitbucket.org:pmcgrath/bbdotfiles.git", "git@github.com:pmcgrath/dotfiles.git", "git@github.com:pmcgrath/other.git"}
	remoteName := "some_other_remote"

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
			t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoUrl, repoResult.RepoPath)
		}
		if repoResult.Command != cloneCmd {
			t.Errorf("Unexpected result Command, expected %s, but got %s", cloneCmd, repoResult.Command)
		}
		// Sample command is : git clone --origin github git@github.com:pmcgrath/dotfiles.git
		expectedOutput := fmt.Sprintf("cmdName: git, args: clone --origin %s %s", remoteName, repoUrl)
		if repoResult.Output[0] != expectedOutput {
			t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
		}
	}
}

func TestExecGitFetch(t *testing.T) {
	// Needed to be in sorted order as the test assertions depend on this
	repoPaths := make([]string, 100)
	for index := 0; index < len(repoPaths); index++ {
		repoPaths[index] = fmt.Sprintf("/tmp/repos/repo%03d", index+1)
	}
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
			t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
		}
		if repoResult.Command != fetchCmd {
			t.Errorf("Unexpected result Command, expected %s, but got %s", fetchCmd, repoResult.Command)

		}
		// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 fetch origin
		expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s fetch %s", repoPath, repoPath, remoteName)
		if repoResult.Output[0] != expectedOutput {
			t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
		}
	}
}

func TestExecGitPull(t *testing.T) {
	repoPath := "/tmp/repos/repo1"
	remoteName := "upstream"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitPull([]string{repoPath}, remoteName)

	if len(results) != 1 {
		t.Errorf("Unexpected result count : %d", len(results))
	}

	repoResult := results[0]

	if repoResult.RepoPath != repoPath {
		t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
	}
	if repoResult.Command != pullCmd {
		t.Errorf("Unexpected result Command, expected %s, but got %s", pullCmd, repoResult.Command)

	}
	// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 pull origin
	expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s pull %s", repoPath, repoPath, remoteName)
	if repoResult.Output[0] != expectedOutput {
		t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
	}
}

func TestExecGitRemote(t *testing.T) {
	repoPath := "/tmp/repos/repo1"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitRemote([]string{repoPath})

	if len(results) != 1 {
		t.Errorf("Unexpected result count : %d", len(results))
	}

	repoResult := results[0]

	if repoResult.RepoPath != repoPath {
		t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
	}
	if repoResult.Command != remoteCmd {
		t.Errorf("Unexpected result Command, expected %s, but got %s", remoteCmd, repoResult.Command)

	}
	// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 remote --verbose
	expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s remote --verbose", repoPath, repoPath)
	if repoResult.Output[0] != expectedOutput {
		t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
	}
}

func TestExecGitStatus(t *testing.T) {
	repoPath := "/tmp/repos/repo1"

	previous := cmdExecutionFn
	defer func() { cmdExecutionFn = previous }()
	cmdExecutionFn = execCmdSpy

	results := execGitStatus([]string{repoPath})

	if len(results) != 1 {
		t.Errorf("Unexpected result count : %d", len(results))
	}

	repoResult := results[0]

	if repoResult.RepoPath != repoPath {
		t.Errorf("Unexpected result RepoPath, expected %s, but got %s", repoPath, repoResult.RepoPath)
	}
	if repoResult.Command != statusCmd {
		t.Errorf("Unexpected result Command, expected %s, but got %s", statusCmd, repoResult.Command)

	}
	// Sample command is : git --git-dir /tmp/repos/repo1/.git --work-tree /tmp/repos/repo1 status --porcelain
	expectedOutput := fmt.Sprintf("cmdName: git, args: --git-dir %s/.git --work-tree %s status --porcelain", repoPath, repoPath)
	if repoResult.Output[0] != expectedOutput {
		t.Errorf("Unexpected result Output, expected %s, but got %s", expectedOutput, repoResult.Output[0])
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
