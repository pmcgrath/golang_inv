package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	searchDir := "c:/repos/stash/ser"

	directoryNames, _ := getSubDirectoryNames(searchDir)
	for _, directoryName := range directoryNames {
		fmt.Println(directoryName)
	}
}

func getSubDirectoryNames(directoryPath string) (names []string, err error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return
	}

	for _, candidate := range files {
		if candidate.IsDir() {
			names = append(names, candidate.Name())
		}
	}

	return
}
