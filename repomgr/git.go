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

func execGitBranch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "branch", "-av")
}

func execGitClone(rootDirectoryPath string, repoUrls []string, remoteName string) gitCmdResults {
	os.Chdir(rootDirectoryPath)

	// This executes clone on each of the url's - assumes the repo is not already cloned
	repoCount := len(repoUrls)
	resultsCh := make(chan gitCmdResult, repoCount)

	var wg sync.WaitGroup
	for _, repoUrl := range repoUrls {
		wg.Add(1)
		go func(repoUrl string) {
			defer wg.Done()

			logDebugf("About to clone for [%s]\n", repoUrl)
			cmdArgs := []string{"clone", "--origin", remoteName, repoUrl}
			cmdOutput, err := execCmd("git", cmdArgs...)
			resultsCh <- gitCmdResult{RepoPath: repoUrl, Command: "clone", Output: cmdOutput, Error: err}
			logDebugf("Completed clone for [%s]\n", repoUrl)
		}(repoUrl)
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
func execGitFetch(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "fetch", remoteName)
}

func execGitPull(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "pull", remoteName)
}

func execGitRemote(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "remote", "--verbose")
}

func execGitStatus(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "status", "--porcelain")
}

func getGitRepoNameFromUrl(repoUrl string) string {
	_, repoDirectoryName := path.Split(repoUrl)
	return strings.TrimSuffix(repoDirectoryName, ".git")
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

func execGitCmdOnMultipleRepos(repoPaths []string, command string, args ...string) gitCmdResults {
	// This executes the same command on a bunch of already existing repos
	repoCount := len(repoPaths)
	resultsCh := make(chan gitCmdResult, repoCount)

	var wg sync.WaitGroup
	for _, repoPath := range repoPaths {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()

			// Need to use --git-dir and --work-tree git args, was using a os.Chdir golang instruction but this was changing the working dir for the
			// golang process and then trying to run a git command, but since we are using goroutines this is unpredictable, where a number of them
			// can be changing the dir at the same time, could end up running a git command in a different directory, by using these git args, the
			// process can control this for each repo. We are assuming the .git directory is a sub directory within the git repo which is the default
			gitDir := path.Join(repoPath, ".git")
			gitArgs := append([]string{"--git-dir", gitDir, "--work-tree", repoPath, command}, args...)
			logDebugf("About to run git with the following args %v", gitArgs)
			cmdOutput, err := execCmd("git", gitArgs...)
			resultsCh <- gitCmdResult{RepoPath: repoPath, Command: command, Output: cmdOutput, Error: err}
		}(repoPath)
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
