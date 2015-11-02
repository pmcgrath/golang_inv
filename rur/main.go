package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	projectsDirectoryPath := flag.String("projectsdirectorypath", "c:/repos/stash/ser", "Projects directory path")

	subDirectoryPaths, err := getAllSubDirectoryPaths(*projectsDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}

	var svcConfigs []serviceConfiguration
	for _, subDirectoryPath := range subDirectoryPaths {
		if !testIfFileExists(subDirectoryPath + "/dn-ci-runner.ps1") {
			continue
		}

		svcConfig, err := parseForService(subDirectoryPath)
		if err != nil {
			log.Fatal(err)
		}
		svcConfigs = append(svcConfigs, svcConfig)

		fmt.Printf("\n\n*******\n")
		fmt.Println(svcConfig)
	}

	for _, svcConfig := range svcConfigs {
		fmt.Printf("\n\n**^^^^^^^^^*****\n%s", svcConfig.Name)
		if svcConfig.Name == "service1" {
			fmt.Printf("\n\n*******\n%s", svcConfig)
		}
	}
	/*
		envResourceMaps := make(map[string]map[string][]string)
		for _, svcConfig := range svcConfigs {
			for env, envConfig := range svcConfig.Environments {
				envResourceMap, ok := envResourceMaps[env]
				if !ok {
					envResourceMap = make(map[string][]string)
				}
				for _, database := range envConfig.Databases {
					envResourceMap[svcConfig.Name] = append(envResourceMap[svcConfig.Name], fmt.Sprintf("DB Type=%s, Host=%s, Name=%s", database.Type, database.Host, database.Name))
				}
				for _, logger := range envConfig.Loggers {
					envResourceMap[svcConfig.Name] = append(envResourceMap[svcConfig.Name], fmt.Sprintf("LOGGER Destination=%s", logger.Destination))
				}

				envResourceMaps[env] = envResourceMap
			}
		}

		for env, envResourceMap := range envResourceMaps {
			fmt.Printf("\n\nENV: %s\n", env)

			for svc, resource := range envResourceMap {
				for _, resourceItem := range resource {
					fmt.Printf("\t%s --> %s\n", svc, resourceItem)
				}
			}
		}
	*/
}
