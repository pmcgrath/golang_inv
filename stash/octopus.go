/*
	See

		http://octopus/api/projects/all
		http://octopus/api/variables/variableset-Projects-19
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type octopus struct {
	connAttrs   providerConnectionAttributes
	timeoutInMS int
}

type octopusProject struct {
	name           string
	id             string
	progressionUrl string
	variablesUrl   string
}

type octopusProjects []octopusProject

func newOctopusProvider(connAttrs providerConnectionAttributes) octopus {
	return octopus{connAttrs: connAttrs, timeoutInMS: 15000}
}

func (p octopus) getHttpClient() http.Client {
	timeout := time.Duration(time.Duration(p.timeoutInMS) * time.Millisecond)

	return http.Client{
		Timeout: timeout,
	}
}

func (p octopus) getProjects() (projects octopusProjects, err error) {
	client := p.getHttpClient()

	url := fmt.Sprintf("%s/api/projects/all", p.connAttrs.URL)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Octopus-ApiKey", p.connAttrs.APIKey)
	req.Header.Set("Accept", "application/json")

	//logDebugf("About to make a http get on [%s]\n", url)
	log.Printf("********************* About to make a http get on [%s]\n", url)
	resp, err := client.Do(req)
	if resp != nil {
		// See http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non 200 status code for [%s], code was %d", url, resp.StatusCode)
		return
	}

	var respData []struct {
		Id    string
		Name  string
		Links map[string]string
	}

	//logDebugln("About to decode response data")
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return
	}

	for _, respDataValue := range respData {
		project := octopusProject{
			name: respDataValue.Name,
			id:   respDataValue.Id,
		}

		project.progressionUrl, _ = respDataValue.Links["Progression"]
		project.variablesUrl, _ = respDataValue.Links["Variables"]

		projects = append(projects, project)
	}

	return
}
