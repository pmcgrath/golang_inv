package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

type gitCmdResult struct {
	RepoPath string
	Command  string
	Output   []string
	Error    error
}

type gitCmdResults []gitCmdResult

func (r gitCmdResults) Len() int {
	return len(r)
}

func (r gitCmdResults) Less(i, j int) bool {
	return r[i].RepoPath < r[j].RepoPath
}

func (r gitCmdResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
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

func execGitClone(rootDirectoryPath string, repoUrls []string, args ...string) gitCmdResults {
	// This executes clone on each of the url's - assumes the repo is not already cloned
	repoCount := len(repoUrls)
	resultsCh := make(chan gitCmdResult, repoCount)

	var wg sync.WaitGroup
	for _, repoUrl := range repoUrls {
		wg.Add(1)
		go func(repoUrl string) {
			defer wg.Done()

			cmdArgs := append([]string{"clone", repoUrl}, args...)
			cmdOutput, err := execCmd("git", cmdArgs...) // BROKE !!
			resultsCh <- gitCmdResult{RepoPath: repoUrl, Command: "clone", Output: cmdOutput, Error: err}
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

func execGitBranch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "symbolic-ref", "--short", "-q", "HEAD")
}

func execGitFetch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "fetch", "origin")
}

func execGitPull(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "pull", "origin")
}

func execGitStatus(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "status", "--porcelain")
}

func getGitRepoNameFromUrl(repoUrl string) string {
	_, repoDirectoryName := path.Split(repoUrl)
	return strings.TrimSuffix(repoDirectoryName, ".git")
}

/*
	Should these go
*/
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

func displayGitCmdResults(results gitCmdResults) {
	for _, result := range results {
		fmt.Printf("%s\n", result.RepoPath)
		if result.Error != nil {
			fmt.Printf("\x1b[31m%s\x1b[39;49m\n", result.Error)
		}
		for _, line := range result.Output {
			fmt.Printf("\t%s\n", line)
		}
		fmt.Println()
	}
}

func timeAndDisplay(context string, repoPaths []string, opFunc func([]string) gitCmdResults) {
	fmt.Printf("\n\n*** %s Starting\n", context)
	start := time.Now()
	displayGitCmdResults(opFunc(repoPaths))
	elapsed := time.Since(start)
	fmt.Printf("*** %s Completed in %s\n", context, elapsed)
}

func run() {
	action := "list"
	rootDirectoryPath := "/tmp/repos"
	fileName := "ted"

	fmt.Println("Context\n=======\n")
	fmt.Println(action)
	fmt.Println(rootDirectoryPath)
	fmt.Println(fileName)
	fmt.Println("=======\n")

	directoryPaths, _ := getAllSubDirectoryPaths(rootDirectoryPath)
	repoPaths := filterGitReposOnly(directoryPaths)

	timeAndDisplay("Branch", repoPaths, execGitBranch)
	timeAndDisplay("Status", repoPaths, execGitStatus)
	timeAndDisplay("Fetch", repoPaths, execGitFetch)
	timeAndDisplay("Pull", repoPaths, execGitPull)
}
