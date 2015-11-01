package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	projectsDirectoryPath := flag.String("projectsdirectorypath", "/tmp/services", "Projects directory path")

	subDirectoryPaths, err := getAllSubDirectoryPaths(*projectsDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, subDirectoryPath := range subDirectoryPaths {
		if !testIfFileExists(subDirectoryPath + "/dn-ci-runner.ps1") {
			continue
		}

		configuration, err := parseForService(subDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n\n*******\n")
		fmt.Println(configuration)
	}
}
