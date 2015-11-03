package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	queryOctopus()
	//queryStash()
}

func queryOctopus() {
	url := "http://octopus"
	apiKey := os.Getenv("API_KEY")

	connAttrs := providerConnectionAttributes{
		URL:    url,
		APIKey: apiKey,
	}
	provider := newOctopusProvider(connAttrs)

	log.Println("About to get projects for provider")
	projects, err := provider.getProjects()
	checkError(err)

	for _, project := range projects {
		log.Printf("\n\n\nProcessing project: %#v", project)
	}
}

func queryStash() {
	url := "http://stash"
	parentName := "ser"
	userName := "pmcgrath"
	password := os.Getenv("REPO_PASSWORD")

	connAttrs := providerConnectionAttributes{
		URL:      url,
		Username: userName,
		Password: password,
	}
	provider := newStashProvider(connAttrs)

	log.Println("About to get repositories for provider")
	repos, err := provider.getRepositories(parentName)
	checkError(err)

	dnRepos := make(map[string][]string, 0)
	for _, repo := range repos {
		log.Printf("\n\n\nProcessing repository: %s", repo.Name)
		repoFilePaths, err := provider.getRepositoryFilePaths(repo.ParentName, repo.Name, "/")
		checkError(err)

		serviceProjectPathPrefix := repo.Name + "/"
		isADotNetRepo := false
		var repoConfigFilePaths []string
		for _, repoFilePath := range repoFilePaths {
			if repoFilePath == "dn-ci-runner.ps1" {
				isADotNetRepo = true
			} else if strings.HasPrefix(repoFilePath, serviceProjectPathPrefix) && path.Ext(repoFilePath) == ".config" {
				repoConfigFilePaths = append(repoConfigFilePaths, repoFilePath)
			}
		}
		if isADotNetRepo && len(repoConfigFilePaths) > 0 {
			log.Print(" Config files found")
			dnRepos[repo.Name] = repoConfigFilePaths
		}
	}

	for _, repo := range repos {
		log.Printf("\n\n\nProcessing repository: %s", repo.Name)
		if repoConfigFilePaths, ok := dnRepos[repo.Name]; ok {
			log.Printf(" Is a DN repo")
			for _, repoConfigFilePath := range repoConfigFilePaths {
				log.Printf("\t%s", repoConfigFilePath)
				content, err := provider.getRepositoryFileContent(repo.ParentName, repo.Name, repoConfigFilePath)
				checkError(err)

				log.Printf(">>\n%s<<\n", content)
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
