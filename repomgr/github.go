/*
	https://developer.github.com/v3/

	Can see the data we are working with using
		curl -v -H "Accept: application/json" https://api.github.com/users/pmcgrath/repos
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type gitHub struct {
	connAttrs connectionAttributes
}

func newGitHubProvider(connAttrs connectionAttributes) gitHub {
	return gitHub{connAttrs: connAttrs}
}

func (p gitHub) getRepos(parentName string) (repos []repositoryDetail, err error) {
	timeoutInMS := 15000

	timeout := time.Duration(time.Duration(timeoutInMS) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	url := fmt.Sprintf("%s/users/%s/repos", p.connAttrs.Url, parentName)
	req, err := http.NewRequest("GET", url, nil)
	// Github does not require authentication for public repos
	if p.connAttrs.Password != "" {
		req.SetBasicAuth(p.connAttrs.Username, p.connAttrs.Password)
	}
	req.Header.Set("Accept", "application/json")

	logDebugf("About to make github api call on [%s]\n", url)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non 200 status code for [%s], code was %d", url, resp.StatusCode)
	}

	logDebugln("About to decode response data")
	var respData []map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	for _, repoData := range respData {
		repo := repositoryDetail{
			ParentName: parentName,
			Name:       repoData["name"].(string),
			ProtocolUrls: map[string]string{
				"https": repoData["clone_url"].(string),
				"ssh":   repoData["ssh_url"].(string),
			},
		}
		repos = append(repos, repo)
	}

	return
}
