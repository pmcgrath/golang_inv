/*
	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html

	Can see the data we are working with using
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects?start=0&limit=2
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects/ser/repos?start=0&limit=20
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type createStashPagedUrl func(int, int) string

type processStashMap func(map[string]interface{}) error

type stash struct {
	connAttrs providerConnectionAttributes
}

func newStashProvider(connAttrs providerConnectionAttributes) stash {
	return stash{connAttrs: connAttrs}
}

func (p stash) getRepos(parentName string) (repos repositoryDetails, err error) {
	var projectKeys []string
	if parentName == "" {
		logDebugln("About to get project keys")
		if projectKeys, err = p.getProjectKeys(); err != nil {
			return
		}
	} else {
		projectKeys = append(projectKeys, parentName)
	}

	for _, projectKey := range projectKeys {
		logDebugf("About to get projects for project key [%s]\n", projectKey)
		projectRepos, err := p.getProjectRepos(projectKey)
		if err != nil {
			return nil, err
		}

		repos = append(repos, projectRepos...)
	}

	return
}

func (p stash) getProjectKeys() (projectKeys []string, err error) {
	err = p.processPagedData(
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects?start=%d&limit=%d", p.connAttrs.Url, start, limit)
		},
		func(project map[string]interface{}) error {
			projectKey := project["key"].(string)
			projectKeys = append(projectKeys, projectKey)
			return nil
		})

	return
}

func (p stash) getProjectRepos(projectKey string) (repos repositoryDetails, err error) {
	err = p.processPagedData(
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos?start=%d&limit=%d", p.connAttrs.Url, projectKey, start, limit)
		},
		func(repository map[string]interface{}) error {
			logDebugln("About to process project repo data")
			name := repository["name"].(string)
			links := repository["links"].(map[string]interface{})
			cloneLinks := links["clone"].([]interface{})

			protocolUrls := make(map[string]string)
			for _, cloneLink := range cloneLinks {
				cloneLinkMap := cloneLink.(map[string]interface{})
				name := cloneLinkMap["name"].(string)
				href := cloneLinkMap["href"].(string)
				protocolUrls[name] = href
			}

			repo := repositoryDetail{
				ParentName:   projectKey,
				Name:         name,
				ProtocolUrls: protocolUrls,
			}
			repos = append(repos, repo)

			return nil
		})

	return
}

func (p stash) processPagedData(createUrl createStashPagedUrl, processValue processStashMap) (err error) {
	timeoutInMS, start, limit := 15000, 0, 25

	timeout := time.Duration(time.Duration(timeoutInMS) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	for {
		url := createUrl(start, limit)
		req, err := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(p.connAttrs.Username, p.connAttrs.Password)
		req.Header.Set("Accept", "application/json")

		logDebugf("About to make a http get on [%s]\n", url)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Non 200 status code for [%s], code was %d", url, resp.StatusCode)
		}

		logDebugln("About to decode response data")
		var respData struct {
			Size       int
			IsLastPage bool
			Values     []map[string]interface{}
		}
		if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return err
		}

		for _, respDataValue := range respData.Values {
			if err = processValue(respDataValue); err != nil {
				return err
			}
		}

		if respData.IsLastPage {
			break
		}

		start += limit + 1
	}

	return
}
