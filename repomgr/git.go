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

func createGitCmdArgsForCloningARepo(repoUrl string, remoteName string) []string {
	return []string{"clone", "--origin", remoteName, repoUrl}
}

func createGitCmdArgsForExistingRepo(repoPath string, command string, args ...string) []string {
	// Need to use --git-dir and --work-tree git args, was using a os.Chdir golang instruction but this was changing the working dir for the
	// golang process and then trying to run a git command, but since we are using goroutines this is unpredictable, where a number of them
	// can be changing the dir at the same time, could end up running a git command in a different directory, by using these git args, the
	// process can control this for each repo. We are assuming the .git directory is a sub directory within the git repo which is the default
	gitDir := path.Join(repoPath, ".git")
	return append([]string{"--git-dir", gitDir, "--work-tree", repoPath, command}, args...)
}

func execGitCmdOnMultipleRepos(repos []string, createGitCmdArgs func(string) []string) gitCmdResults {
	repoCount := len(repos)
	resultsCh := make(chan gitCmdResult, repoCount)

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()

			cmdArgs := createGitCmdArgs(repo)
			logDebugf("About to run git with the following args %v", cmdArgs)
			cmdOutput, err := execCmd("git", cmdArgs...)
			logDebugf("Completed running git with the following args %v", cmdArgs)
			resultsCh <- gitCmdResult{Repo: repo, Command: "command", Output: cmdOutput, Error: err}
		}(repo)
	}
	wg.Wait()
	close(resultsCh)

	var res gitCmdResults
	for result := range resultsCh {
		logDebugf("----> %#v\n", result)
		res = append(res, result)
	}

	sort.Sort(res)
	return res
}

func execGitBranch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths,
		func(repoPath string) []string {
			return createGitCmdArgsForExistingRepo(repoPath, "branch", "symbolic-ref", "--short", "-q", "HEAD")
		})
}

func execGitClone(rootDirectoryPath string, repoUrls []string, remoteName string) gitCmdResults {
	// We need to change into the root dir so we can clone into this location
	os.Chdir(rootDirectoryPath)

	// This executes clone on each of the url's - assumes the repo is not already cloned
	return execGitCmdOnMultipleRepos(repoUrls,
		func(repoUrl string) []string {
			return createGitCmdArgsForCloningARepo(repoUrl, remoteName)
		})
}

func execGitFetch(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths,
		func(repoPath string) []string {
			return createGitCmdArgsForExistingRepo(repoPath, "fetch", remoteName)
		})
}

func execGitPull(repoPaths []string, remoteName string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths,
		func(repoPath string) []string {
			return createGitCmdArgsForExistingRepo(repoPath, "pull", remoteName)
		})
}

func execGitStatus(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths,
		func(repoPath string) []string {
			return createGitCmdArgsForExistingRepo(repoPath, "status", "--porcelain")
		})
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
			res = append(res, result.Repo)
		}
	}

	return res
}
