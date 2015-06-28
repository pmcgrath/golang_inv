package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

type gitCmdResult struct {
	RepoPath string
	Command  string
	Result   []string
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

var (
	action            string
	rootDirectoryPath string
	fileName          string
)

func init() {
	initAction := flag.String("action", "status", "Action to run")
	//initRootDirectoryPath := flag.String("rootDirectoryPath", "/home/pmcgrath/go/src/github.com/pmcgrath", "Root directory path")
	initRootDirectoryPath := flag.String("rootDirectoryPath", "/tmp/gitinv", "Root directory path")
	initFileName := flag.String("fileName", "", "File name")
	flag.Parse()

	action = *initAction
	rootDirectoryPath = *initRootDirectoryPath
	fileName = *initFileName
}

func getAllSubDirs(directoryPath string) ([]string, error) {
	var res []string

	files, err := ioutil.ReadDir(rootDirectoryPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsDir() {
			directoryPath := path.Join(directoryPath, file.Name())
			res = append(res, directoryPath)
		}
	}

	return res, nil
}

func execGitCmd(repoPath string, args ...string) ([]string, error) {
	if err := os.Chdir(repoPath); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", args...)
	cmdOutputBytes, err := cmd.CombinedOutput()
	cmdOutputString := string(cmdOutputBytes)
	if strings.TrimSpace(cmdOutputString) == "" {
		// No output (No error) or a raw error
		return nil, err
	}

	if err != nil {
		// Use output as error
		return nil, fmt.Errorf(cmdOutputString)
	}

	outputLines := strings.Split(cmdOutputString, "\n")

	return outputLines, err
}

func execGitCmdOnMultipleRepos(repoPaths []string, command string, args ...string) gitCmdResults {
	repoCount := len(repoPaths)
	resultsCh := make(chan gitCmdResult, repoCount)

	args = append([]string{command}, args...)

	var wg sync.WaitGroup
	for _, repoPath := range repoPaths {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()
			cmdResult, err := execGitCmd(repoPath, args...)
			resultsCh <- gitCmdResult{RepoPath: repoPath, Command: command, Result: cmdResult, Error: err}
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

func getBranch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "symbolic-ref", "--short", "-q", "HEAD")
}

func fetch(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "fetch", "origin")
}

func getStatus(repoPaths []string) gitCmdResults {
	return execGitCmdOnMultipleRepos(repoPaths, "status", "--porcelain")
}

func display(results gitCmdResults) {
	for _, result := range results {
		fmt.Printf("%s\n", result.RepoPath)
		if result.Error != nil {
			fmt.Printf("\x1b[31m%s\x1b[39;49m\n", result.Error)
		}
		for _, line := range result.Result {
			fmt.Printf("\t%s\n", line)
		}
		fmt.Println()
	}
}

func filterGitReposOnly(directoryPaths []string) []string {
	var res []string

	results := getBranch(directoryPaths)
	for _, result := range results {
		if result.Error == nil {
			res = append(res, result.RepoPath)
		}
	}

	return res
}

func timeAndDisplay(context string, repoPaths []string, opFunc func([]string) gitCmdResults) {
	start := time.Now()
	display(opFunc(repoPaths))
	elapsed := time.Since(start)
	fmt.Printf("*** %s elapsed is %s\n", context, elapsed)
}

func main() {
	fmt.Println(action)
	fmt.Println(rootDirectoryPath)
	fmt.Println(fileName)

	directoryPaths, _ := getAllSubDirs(rootDirectoryPath)
	repoPaths := filterGitReposOnly(directoryPaths)

	timeAndDisplay("Branch", repoPaths, getBranch)
	timeAndDisplay("Fetch", repoPaths, fetch)
	timeAndDisplay("Status", repoPaths, getStatus)
}
