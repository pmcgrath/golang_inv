package main

import (
	"io/ioutil"
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

func getAllSubDirs(directoryPath string) ([]string, error) {
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
