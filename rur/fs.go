package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getAllSubDirectoryPaths(directoryPath string) ([]string, error) {
	dirs := []string{}

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

func getDirectoryFiles(directoryPath, pattern string) ([]string, error) {
	globPattern := path.Join(directoryPath, pattern)

	files := []string{}
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		// filepath.Glob function return Windows style seperators "\", but we use path.Split which does not seem to work so switch path seperators back to *nix tyle
		match = strings.Replace(match, "\\", "/", -1)

		matchFileInfo, err := os.Stat(match)
		if err != nil {
			return nil, err
		}

		if !matchFileInfo.Mode().IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

func testIfFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}
