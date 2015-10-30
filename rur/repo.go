package main

import (
	"fmt"
	"path"
	"strings"
)

func getFSBasedRepoConfigFilePaths(directoryPath string) (configFilePaths []string, err error) {
	serviceName := path.Base(directoryPath)

	mainProjectDirectoryPath := path.Join(directoryPath, serviceName)
	allConfigFilePaths, err := getDirectoryFiles(mainProjectDirectoryPath, "*.config")
	if err != nil {
		return
	}
	if len(allConfigFilePaths) == 0 {
		err = fmt.Errorf("No config file exists for service %s", serviceName)
		return
	}

	for _, configFilePath := range allConfigFilePaths {
		if isCandidateConfigFile(configFilePath) {
			configFilePaths = append(configFilePaths, configFilePath)
		}
	}

	return
}

func isCandidateConfigFile(filePath string) bool {
	// packages.config, *.debug.config and *.release.config files are not candidates
	_, fileName := path.Split(filePath)
	fileName = strings.ToLower(fileName)
	if fileName == "packages.config" || strings.HasSuffix(fileName, ".debug.config") || strings.HasSuffix(fileName, ".release.config") {
		return false
	}

	return true
}
