/*	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getStashRepoDetails(connAttrs connectionAttributes) (result []repoDetail, err error) {
	var projectKeys []string
	if connAttrs.ParentName == "" {
		if projectKeys, err = getStashProjectKeys(connAttrs); err != nil {
			return
		}
	} else {
		projectKeys = append(projectKeys, connAttrs.ParentName)
	}

	for _, projectKey := range projectKeys {
		projectRepos, err := getStashProjectRepos(connAttrs, projectKey)
		if err != nil {
			return nil, err
		}

		result = append(result, projectRepos...)
	}

	return
}

func getStashProjectKeys(connAttrs connectionAttributes) (projectKeys []string, err error) {
	err = processStashPagedData(
		connAttrs,
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects?start=%d&limit=%d", connAttrs.Url, start, limit)
		},
		func(project map[string]interface{}) error {
			projectKey := project["key"].(string)
			projectKeys = append(projectKeys, projectKey)
			return nil
		})

	return
}

func getStashProjectRepos(connAttrs connectionAttributes, projectKey string) (repos []repoDetail, err error) {
	err = processStashPagedData(
		connAttrs,
		func(start, limit int) string {
			return fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos?start=%d&limit=%d", connAttrs.Url, projectKey, start, limit)
		},
		func(repository map[string]interface{}) error {
			name := repository["name"].(string)
			links := repository["links"].(map[string]interface{})
			cloneLinks := links["clone"].([]map[string]string)
			protocolUrls := make(map[string]string)
			for _, cloneLink := range cloneLinks {
				protocolUrls[cloneLink["name"]] = cloneLink["href"]
			}

			repo := repoDetail{
				ParentName:   connAttrs.ParentName,
				Name:         name,
				ProtocolUrls: protocolUrls,
			}
			repos = append(repos, repo)

			return nil
		})

	return
}

func processStashPagedData(connAttrs connectionAttributes, createUrl func(int, int) string, processValue func(map[string]interface{}) error) (err error) {
	timeoutInMS, start, limit := 15000, 0, 25

	timeout := time.Duration(time.Duration(timeoutInMS) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	for {
		stashUrl := createUrl(start, limit)
		req, err := http.NewRequest("GET", stashUrl, nil)
		req.SetBasicAuth(connAttrs.Username, connAttrs.Password)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Non 200 status code for [%s], code was %d", stashUrl, resp.StatusCode)
		}

		var data struct {
			Size       int
			IsLastPage bool
			Values     []map[string]interface{}
		}
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		for _, value := range data.Values {
			if err = processValue(value); err != nil {
				return err
			}
		}

		if data.IsLastPage {
			break
		}

		start += limit + 1
	}

	return
}
