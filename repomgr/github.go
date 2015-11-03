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
	connAttrs providerConnectionAttributes
}

func newGitHubProvider(connAttrs providerConnectionAttributes) gitHub {
	return gitHub{connAttrs: connAttrs}
}

func (p gitHub) getRepos(parentName string) (repos repositoryDetails, err error) {
	timeoutInMS := 15000

	timeout := time.Duration(time.Duration(timeoutInMS) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	url := fmt.Sprintf("%s/users/%s/repos", p.connAttrs.URL, parentName)
	req, err := http.NewRequest("GET", url, nil)
	// Github does not require authentication for public repos
	if p.connAttrs.Password != "" {
		req.SetBasicAuth(p.connAttrs.Username, p.connAttrs.Password)
	}
	req.Header.Set("Accept", "application/json")

	logDebugf("About to make github api call on [%s]\n", url)
	resp, err := client.Do(req)
	if resp != nil {
		// See http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

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
			ProtocolURLs: map[string]string{
				"https": repoData["clone_url"].(string),
				"ssh":   repoData["ssh_url"].(string),
			},
		}
		repos = append(repos, repo)
	}

	return
}
