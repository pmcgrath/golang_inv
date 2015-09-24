/*
	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Stash content
type stashConfiguration interface {
	getUrl() string
	getUsername() string
	getPassword() string
}

func (c configuration) getUrl() string {
	return c.stashUrl
}

func (c configuration) getUsername() string {
	return c.stashUsername
}

func (c configuration) getPassword() string {
	return c.stashPassword
}

func showStashProjectListing(config stashConfiguration, projectKey string) {
	projectKeys, err := getStashProjectKeys(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting")
	for _, candidateProjectKey := range projectKeys {
		if projectKey != "" && candidateProjectKey != projectKey {
			fmt.Printf("Checking %s  against %s\n", candidateProjectKey, projectKey)
			continue
		}
		fmt.Printf("Project key : %s\n", candidateProjectKey)

		projectRepos, err := getStashProjectRepos(config, candidateProjectKey)
		if err != nil {
			log.Fatal(err)
		}

		for _, projectRepo := range projectRepos {
			fmt.Printf("\tRepo : %s\n", projectRepo)
		}
	}
}

func getStashProjectRepos(config stashConfiguration, projectKey string) (projectRepos []string, err error) {
	return getListFromStashPagedData(
		config,
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos?start=%d&limit=%d", config.getUrl(), projectKey, start, limit)
		},
		"name")
}

func getStashProjectKeys(config stashConfiguration) (projectKeys []string, err error) {
	return getListFromStashPagedData(
		config,
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects?start=%d&limit=%d", config.getUrl(), start, limit)
		},
		"key")
}

func getListFromStashPagedData(config stashConfiguration, createUrl func(int, int) string, key string) (list []string, err error) {
	timeoutInMS, start, limit := 15000, 0, 25

	timeout := time.Duration(time.Duration(timeoutInMS) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	for {
		stashUrl := createUrl(start, limit)

		req, err := http.NewRequest("GET", stashUrl, nil)
		req.SetBasicAuth(config.getUsername(), config.getPassword())
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var data struct {
			Size       int
			IsLastPage bool
			Values     []map[string]interface{}
		}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		for _, project := range data.Values {
			list = append(list, project[key].(string))
		}

		if data.IsLastPage {
			break
		}

		start += limit + 1
	}

	return
}
