package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	connectionString := "Server=tcp:dbserver1; Database=myDB; MultiSubnetFailover=True; Integrated Security=SSPI;"
	parseMsSqlConnectionString(connectionString)

	projectsDirectoryPath := flag.String("projectsdirectorypath", "c:/repos/stash/ser", "Projects directory path")

	subDirectoryPaths, err := getAllSubDirectoryPaths(*projectsDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, subDirectoryPath := range subDirectoryPaths {
		configuration, err := parseForService(subDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n\n*******\n")
		fmt.Println(configuration)
		return
	}
}
