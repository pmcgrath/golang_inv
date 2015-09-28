package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"sync"
)

func getAllSubDirectoryPaths(directoryPath string) ([]string, error) {
	var dirs []string

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsDir() {
			directoryPath := path.Join(directoryPath, file.Name())
			dirs = append(dirs, directoryPath)
		}
	}

	return dirs, nil
}

type cmdResult struct {
	Path   string
	Result []string
	Error  error
}

type cmdResults []cmdResult

func (r cmdResults) Len() int {
	return len(r)
}

func (r cmdResults) Less(i, j int) bool {
	return r[i].Path < r[j].Path
}

func (r cmdResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func execCmd(workingDirectoryPath string) ([]string, error) {
	if err := os.Chdir(workingDirectoryPath); err != nil {
		return nil, err
	}

	cmd := exec.Command("ls", "-al")
	// PENDING - THIS WILL NOT WOEK IS COMBINING ACROSS goroutines
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

func execCmdOnMultipleDirs(directoryPaths []string) cmdResults {
	// This executes the same command on a bunch of already existing repos
	count := len(directoryPaths)
	resultsCh := make(chan cmdResult, count)

	var wg sync.WaitGroup
	for _, directoryPath := range directoryPaths {
		wg.Add(1)
		go func(directoryPath string) {
			defer wg.Done()
			result, err := execCmd(directoryPath)
			resultsCh <- cmdResult{Path: directoryPath, Result: result, Error: err}
		}(directoryPath)
	}
	wg.Wait()
	close(resultsCh)

	var res cmdResults
	for result := range resultsCh {
		res = append(res, result)
	}

	sort.Sort(res)
	return res
}

func main() {
	//directoryPaths, _ := getAllSubDirectoryPaths("/tmp")
	//directoryPaths := []string{"/tmp/does-not-exist", "/tmp/repos", "/tmp/.com.google.Chrome.XNuGDt/"}
	directoryPaths := []string{"/tmp/.com.google.Chrome.XNuGDt/"}
	results := execCmdOnMultipleDirs(directoryPaths)

	for _, result := range results {
		fmt.Printf("%s\n", result.Path)
		for _, entry := range result.Result {
			fmt.Printf("%#v\n\n", entry)
		}
		fmt.Printf("ERR found: %t --> %#v\n\n", result.Error != nil, result.Error)
	}
}
