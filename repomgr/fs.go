package main

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

func getDefaultProjectsDirectoryPath() string {
	if runtime.GOOS != "windows" {
		dir, _ := getCurrentUserHomedir()
		dir += "/repos"
		return dir
	}

	return "c:/repos"
}

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

func testIfDirectoryExists(directoryPath string) bool {
	info, err := os.Stat(directoryPath)
	return err == nil && info.IsDir()
}
