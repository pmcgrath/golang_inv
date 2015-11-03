/*
	See 	https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp573712
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1574208
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1037296
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1055056
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1055056
		https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp849632

	Can see the data we are working with using
		# List of projects
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects?start=0&limit=2

		# List of repositories for a project (Project key = ser)
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects/ser/repos?start=0&limit=20

		# List of children (directories and files) for a repository path (Project key = ser, repository slug = project1 and path = / )
		# Note the path is case sensitive
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects/ser/repos/project1/browse/?start=0&limit=20

		# List of children (directories and files) for a repository path (Project key = ser, repository slug = project1 and path = /packages )
		# Note the path is case sensitive
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: application/json" https://stash/rest/api/1.0/projects/ser/repos/project1/browse/packages?start=0&limit=20

		# Get raw content, doing so on tip
		curl -v -u "pmcgrath:PASSWORD" -H "Accept: text/plain" https://stash/rest/projects/ser/repos/project1/file1.sh?raw
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type createStashPagedURL func(int, int) string

type processResponseData func() (bool, error)

type stash struct {
	connAttrs   providerConnectionAttributes
	timeoutInMS int
}

func newStashProvider(connAttrs providerConnectionAttributes) stash {
	return stash{connAttrs: connAttrs, timeoutInMS: 15000}
}

func (p stash) getHttpClient() http.Client {
	timeout := time.Duration(time.Duration(p.timeoutInMS) * time.Millisecond)

	return http.Client{
		Timeout: timeout,
	}
}

func (p stash) getRepositories(parentName string) (repos repositoryDetails, err error) {
	var projectKeys []string
	if parentName == "" {
		if projectKeys, err = p.getProjectKeys(); err != nil {
			return
		}
	} else {
		projectKeys = append(projectKeys, parentName)
	}

	for _, projectKey := range projectKeys {
		projectRepos, err := p.getProjectRepositories(projectKey)
		if err != nil {
			return nil, err
		}

		repos = append(repos, projectRepos...)
	}

	return
}

func (p stash) getProjectKeys() (projectKeys []string, err error) {
	// See http://stash/rest/api/1.0/projects
	var respData struct {
		Size       int
		IsLastPage bool
		Values     []map[string]interface{}
	}

	err = p.processPagedData(
		func(start, limit int) string {
			// See https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp573712
			return fmt.Sprintf("%s/rest/api/1.0/projects?start=%d&limit=%d", p.connAttrs.URL, start, limit)
		},
		&respData,
		func() (bool, error) {
			for _, respDataValue := range respData.Values {
				projectKey := respDataValue["key"].(string)
				projectKeys = append(projectKeys, projectKey)
			}

			return respData.IsLastPage, nil
		})

	return
}

func (p stash) getProjectRepositories(projectKey string) (repos repositoryDetails, err error) {
	// See http://stash/rest/api/1.0/projects/ser/repos
	var respData struct {
		Size       int
		IsLastPage bool
		Values     []map[string]interface{}
	}

	err = p.processPagedData(
		func(start, limit int) string {
			// See https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1574208
			return fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos?start=%d&limit=%d", p.connAttrs.URL, projectKey, start, limit)
		},
		&respData,
		func() (bool, error) {
			for _, respDataValue := range respData.Values {
				name := respDataValue["name"].(string)
				links := respDataValue["links"].(map[string]interface{})
				cloneLinks := links["clone"].([]interface{})

				protocolURLs := make(map[string]string)
				for _, cloneLink := range cloneLinks {
					cloneLinkMap := cloneLink.(map[string]interface{})
					name := cloneLinkMap["name"].(string)
					href := cloneLinkMap["href"].(string)
					protocolURLs[name] = href
				}

				repo := repositoryDetail{
					ParentName:   projectKey,
					Name:         name,
					ProtocolURLs: protocolURLs,
				}
				repos = append(repos, repo)
			}

			return respData.IsLastPage, nil
		})

	return
}

func (p stash) getRepositoryFilePaths(projectKey, repositorySlug, path string) (filePaths []string, err error) {
	// See http://stash/rest/api/1.0/projects/ser/repos/repo1/files?limit=2000
	// Can see sub dir's content with http://stash/rest/api/1.0/projects/ser/repos/repo1/files/packages?limit=2000
	// Seems to be case sensitive for the paths
	var respData struct {
		Size       int
		IsLastPage bool
		Values     []string
	}

	err = p.processPagedData(
		func(start, limit int) string {
			// https://developer.atlassian.com/static/rest/stash/3.11.3/stash-rest.html#idp1037296
			return fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/files%s?start=%d&limit=%d", p.connAttrs.URL, projectKey, repositorySlug, path, start, limit)
		},
		&respData,
		func() (bool, error) {
			filePaths = append(filePaths, respData.Values...)

			return respData.IsLastPage, nil
		})

	return
}

func (p stash) getRepositoryFileContent(projectKey, repositorySlug, path string) (content string, err error) {
	client := p.getHttpClient()

	// Will just use raw get here
	url := fmt.Sprintf("%s/projects/%s/repos/%s/browse/%s?raw", p.connAttrs.URL, projectKey, repositorySlug, path)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(p.connAttrs.Username, p.connAttrs.Password)
	req.Header.Set("Accept", "plain/text")

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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	content = string(bodyBytes)

	return
}

func (p stash) processPagedData(createURL createStashPagedURL, respData interface{}, processData processResponseData) (err error) {
	start, limit := 0, 5000 // Pretty much avoid paging here by using 5000 as limit

	client := p.getHttpClient()

	for {
		url := createURL(start, limit)
		req, err := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(p.connAttrs.Username, p.connAttrs.Password)
		req.Header.Set("Accept", "application/json")

		//logDebugf("About to make a http get on [%s]\n", url)
		log.Printf("********************* About to make a http get on [%s]\n", url)
		resp, err := client.Do(req)
		if resp != nil {
			// See http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body
			defer resp.Body.Close()
		}
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Non 200 status code for [%s], code was %d", url, resp.StatusCode)
		}

		//logDebugln("About to decode response data")
		if err = json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return err
		}

		isLastPage, err := processData()
		if err != nil {
			return err
		}
		if isLastPage {
			break
		}

		start += limit + 1
	}

	return
}
