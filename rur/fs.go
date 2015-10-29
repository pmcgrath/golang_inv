package main

import (
	"io/ioutil"
	"os"
	"path"
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

func testIfFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}
