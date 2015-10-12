/*
	Just shelling out to git here, could have used https://github.com/libgit2/git2go
	Didn't bother with a pool here as we do not expect the repo count to be too big
*/
package main

import (
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

var (
	cmdExecutionFn = execCmd
)

func execGitBranch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleExistingRepos(repoPaths, "branch", "-av")
}

func execGitClone(rootDirectoryPath string, repoUrls []string, remoteName string) gitCmdResults {
	startingDirectoryPath, _ := os.Getwd()
	os.Chdir(rootDirectoryPath)
	defer os.Chdir(startingDirectoryPath)

	// This executes clone on each of the url's - assumes the repo is not already cloned
	return execGitCmdOnMultipleRepos(
		repoUrls,
		"clone",
		func(repoUrl string) []string {
			return []string{"clone", "--origin", remoteName, repoUrl}
		})
}
func execGitFetch(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleExistingRepos(repoPaths, "fetch", remoteName)
}

func execGitPull(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleExistingRepos(repoPaths, "pull", remoteName)
}

func execGitRemote(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleExistingRepos(repoPaths, "remote", "--verbose")
}

func execGitStatus(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleExistingRepos(repoPaths, "status", "--porcelain")
}

func filterGitReposOnly(directoryPaths []string) []string {
	var res []string

	results := execGitBranch(directoryPaths)
	for _, result := range results {
		if result.Error == nil {
			res = append(res, result.RepoPath)
		}
	}

	return res
}

func getGitRepoNameFromUrl(repoUrl string) string {
	_, repoDirectoryName := path.Split(repoUrl)
	return strings.TrimSuffix(repoDirectoryName, ".git")
}

func execGitCmdOnMultipleExistingRepos(repoPaths []string, command string, args ...string) gitCmdResults {
	return execGitCmdOnMultipleRepos(
		repoPaths,
		command,
		func(repoPath string) []string {
			// Need to use --git-dir and --work-tree git args, was using a os.Chdir golang instruction but this was changing the working dir for the
			// golang process and then trying to run a git command, but since we are using goroutines this is unpredictable, where a number of them
			// can be changing the dir at the same time, could end up running a git command in a different directory, by using these git args, the
			// process can control this for each repo. We are assuming the .git directory is a sub directory within the git repo which is the default
			gitDir := path.Join(repoPath, ".git")
			gitArgs := append([]string{"--git-dir", gitDir, "--work-tree", repoPath, command}, args...)

			return gitArgs
		})
}

func execGitCmdOnMultipleRepos(repos []string, gitCommand string, createGitCmdArgs func(string) []string) gitCmdResults {
	repoCount := len(repos)
	resultsCh := make(chan gitCmdResult, repoCount)

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()

			gitCmdArgs := createGitCmdArgs(repo)
			logDebugf("About to run git with the following args %v", gitCmdArgs)
			cmdOutput, err := cmdExecutionFn("git", gitCmdArgs...)
			logDebugf("Completed running git with the following args %v", gitCmdArgs)
			resultsCh <- gitCmdResult{RepoPath: repo, Command: gitCommand, Output: cmdOutput, Error: err}
		}(repo)
	}
	wg.Wait()
	close(resultsCh)

	var res gitCmdResults
	for result := range resultsCh {
		res = append(res, result)
	}

	sort.Sort(res)
	return res
}
