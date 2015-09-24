/*
	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)


type interface sourceProvider {
	getProjectList()
	getRepos	
}




func main() {
	/*
		directoryNames, _ := getSubDirectoryNames(searchDir)
		for _, directoryName := range directoryNames {
			fmt.Println(directoryName)
		}
	*/
	config, err := getConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting")
	showStashProjectListing(config, config.stashProjectKey)

	projectDirectoryPath := path.Join(config.projectsRootDirectoryPath, config.stashProjectKey)
	directoryNames, err := getSubDirectoryNames(projectDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, directoryName := range directoryNames {
		fmt.Printf("Dir -> %s\n", directoryName)
	}

	stashProjectRepos, err := getStashProjectRepos(config, config.stashProjectKey)
	if err != nil {
		log.Fatal(err)
	}

	originalDirectoryPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(originalDirectoryPath)

	projectDirectoryPath = path.Join(config.projectsRootDirectoryPath, config.stashProjectKey)
	os.Chdir(projectDirectoryPath)

	for _, stashProjectRepo := range stashProjectRepos {
		fmt.Printf("Stash repo -> %s\n", stashProjectRepo)
		stashProjectRepo = strings.ToLower(stashProjectRepo)
		repoDirectoryPath := path.Join(projectDirectoryPath, stashProjectRepo)

		fmt.Printf("->> %s    %s   [%s]", projectDirectoryPath, stashProjectRepo, repoDirectoryPath)

		if _, err := os.Stat(repoDirectoryPath); os.IsNotExist(err) {
			fmt.Printf("***** Need to deploy this one %s\n", repoDirectoryPath)
		}
	}
}

// Entry point content
type configuration struct {
	stashUrl                  string
	stashUsername             string
	stashPassword             string
	stashSshUrl               string
	stashProjectKey           string
	projectsRootDirectoryPath string
	useSshProtocol            bool
}

func getConfiguration() (config configuration, err error) {
	// Stash attributes
	stashUrl := flag.String("serverurl", "http://stash", "Stash url - prefix")
	stashUsername := flag.String("username", "", "Stash username - if not supplied will be the current user's name")
	stashPassword := flag.String("password", "", "Stash password - if not supplied will be try to use the STASHPWD environment variable")
	stashSshUrl := flag.String("sshurl", "ssh://git@stash:7999", "Stash ssh url - prefix")
	stashProjectKey := flag.String("projectkey", "SER", "Stash project key")
	// Local setup attributes
	projectsRootDirectoryPath := flag.String("projectsDirectoryPath", "c:/repos/stash", "Projects root directory path")
	useSshProtocol := flag.Bool("usessh", true, "Use ssh")
	// Parse
	flag.Parse()

	config = configuration{
		stashUrl:                  *stashUrl,
		stashUsername:             *stashUsername,
		stashPassword:             *stashPassword,
		stashSshUrl:               *stashSshUrl,
		stashProjectKey:           *stashProjectKey,
		projectsRootDirectoryPath: *projectsRootDirectoryPath,
		useSshProtocol:            *useSshProtocol,
	}

	// Overrides
	if config.stashUsername == "" {
		user, err := user.Current()
		if err != nil {
			return config, err
		}

		config.stashUsername = user.Username
		domainSperatorIndex := strings.Index(config.stashUsername, "\\")
		if domainSperatorIndex > -1 {
			config.stashUsername = config.stashUsername[domainSperatorIndex+1:]
		}
	}
	if config.stashPassword == "" {
		config.stashPassword = os.Getenv("STASHPWD")
	}

	// Validation
	// PENDING
	fmt.Printf("config is %#v\n", config)

	return
}

// File system content
func getSubDirectoryNames(directoryPath string) (names []string, err error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return
	}

	for _, candidate := range files {
		if candidate.IsDir() {
			names = append(names, candidate.Name())
		}
	}

	return
}

/*
    url := "http://restapi3.apiary.io/notes"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))


# Stash attributes
$stashUrl = 'http://stash'
$stashSshUrl = 'ssh://git@stash:7999'
$stashProject = 'SER'

# Local setup attributes
$projectRootDirectoryPath = 'c:\repos\stash\ser'
$useSshProtocol = $true

# Get credentials
$userName = read-host "Gimme your stash user name, default is $($env:UserName)"
$userName = if ($userName -ne '') { $userName } else { $env:UserName }
$password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR((read-host 'Gimme your stash password' -assecurestring)));
$credentials = "$userName`:$password"

# Get stash projects information
$stashProjectsUrl = "$stashUrl/rest/api/1.0/projects"
$projects = convertfrom-json (curl -u $credentials $stashProjectsUrl)
# Echo info - lots of data here
#$projects.Values
# Echo names
#$projects.Values | select name

# Get stash repos information
$stashReposUrl = "$stashUrl/rest/api/1.0/projects/$stashProject/repos"
$repos = convertfrom-json (curl -u $credentials $stashReposUrl)
# Echo info - lots of data here
#$repos.Values
# Echo names
#$repos.Values | select name

# Clone if not already cloned
pushd
cd $projectRootDirectoryPath
foreach ($repoName in ($repos.Values).name)
{
	$repoName = $repoName.ToLower()
	$repoDirectoryPath = join-path $projectRootDirectoryPath $repoName
	if (! (test-path $repoDirectoryPath))
	{
		if ($useSshProtocol)
		{
			$stashRepoUrl = "$stashSshUrl/$stashProject/$repoName.git".ToLower()
			git clone $stashRepoUrl
		}
		else
		{
			write-host 'Have not taken care of http protocol'
		}
	}
}
popd


<#
Pretty version of http://stash/rest/api/1.0/projects/
Have removed all but the first repo - too much noise
{
    "size":4,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "key":"API",
            "id":104,
            "name":"ApiTeam",
            "public":false,
            "type":"NORMAL",
            "link":{
                "url":"/projects/API",
                "rel":"self"
            },
            "links":{
                "self":[
                    {
                        "href":"http://stash/projects/API"
                    }
                ]
            }
        },
		// More projects .... Removed
    ],
    "start":0
}


Pretty version of http://stash/rest/api/1.0/projects/ser/repos
Have removed all but the first repo - too much noise
{
    "size":17,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "slug":"travelrepublic.adverts.service",
            "id":324,
            "name":"TravelRepublic.Adverts.Service",
            "scmId":"git",
            "state":"AVAILABLE",
            "statusMessage":"Available",
            "forkable":true,
            "project":{
	if  err != nil {
		fmt.Printf("Dir -> %s\n", directoryName)
	}

}

// Entry point content
type configuration struct {
	stashUrl                  string
	stashUsername             string
	stashPassword             string
	stashSshUrl               string
	stashProjectKey           string
	projectsRootDirectoryPath string
	useSshProtocol            bool
}

func getConfiguration() (config configuration, err error) {
	// Stash attributes
	stashUrl := flag.String("serverurl", "http://stash", "Stash url - prefix")
	stashUsername := flag.String("username", "", "Stash username - if not supplied will be the current user's name")
	stashPassword := flag.String("password", "", "Stash password - if not supplied will be try to use the STASHPWD environment variable")
	stashSshUrl := flag.String("sshurl", "ssh://git@stash:7999", "Stash ssh url - prefix")
	stashProjectKey := flag.String("projectkey", "SER", "Stash project key")
	// Local setup attributes
	projectsRootDirectoryPath := flag.String("projectsDirectoryPath", "c:\\repos\\stash", "Projects root directory path")
	useSshProtocol := flag.Bool("usessh", true, "Use ssh")
	// Parse
	flag.Parse()

	config = configuration{
		stashUrl:                  *stashUrl,
		stashUsername:             *stashUsername,
		stashPassword:             *stashPassword,
		stashSshUrl:               *stashSshUrl,
		stashProjectKey:           *stashProjectKey,
		projectsRootDirectoryPath: *projectsRootDirectoryPath,
		useSshProtocol:            *useSshProtocol,
	}

	// Overrides
	if config.stashUsername == "" {
		user, err := user.Current()
		if err != nil {
			return config, err
		}

		config.stashUsername = user.Username
		domainSperatorIndex := strings.Index(config.stashUsername, "\\")
		if domainSperatorIndex > -1 {
			config.stashUsername = config.stashUsername[domainSperatorIndex+1:]
		}
	}
	if config.stashPassword == "" {
		config.stashPassword = os.Getenv("STASHPWD")
	}

	// Validation
	// PENDING
	fmt.Printf("config is %#v\n", config)

	return
}

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

// File system content
func getSubDirectoryNames(directoryPath string) (names []string, err error) {
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return
	}

	for _, candidate := range files {
		if candidate.IsDir() {
			names = append(names, candidate.Name())
		}
	}

	return
}

/*
    url := "http://restapi3.apiary.io/notes"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))


# Stash attributes
$stashUrl = 'http://stash'
$stashSshUrl = 'ssh://git@stash:7999'
$stashProject = 'SER'

# Local setup attributes
$projectRootDirectoryPath = 'c:\repos\stash\ser'
$useSshProtocol = $true

# Get credentials
$userName = read-host "Gimme your stash user name, default is $($env:UserName)"
$userName = if ($userName -ne '') { $userName } else { $env:UserName }
$password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR((read-host 'Gimme your stash password' -assecurestring)));
$credentials = "$userName`:$password"

# Get stash projects information
$stashProjectsUrl = "$stashUrl/rest/api/1.0/projects"
$projects = convertfrom-json (curl -u $credentials $stashProjectsUrl)
# Echo info - lots of data here
#$projects.Values
# Echo names
#$projects.Values | select name

# Get stash repos information
$stashReposUrl = "$stashUrl/rest/api/1.0/projects/$stashProject/repos"
$repos = convertfrom-json (curl -u $credentials $stashReposUrl)
# Echo info - lots of data here
#$repos.Values
# Echo names
#$repos.Values | select name

# Clone if not already cloned
pushd
cd $projectRootDirectoryPath
foreach ($repoName in ($repos.Values).name)
{
	$repoName = $repoName.ToLower()
	$repoDirectoryPath = join-path $projectRootDirectoryPath $repoName
	if (! (test-path $repoDirectoryPath))
	{
		if ($useSshProtocol)
		{
			$stashRepoUrl = "$stashSshUrl/$stashProject/$repoName.git".ToLower()
			git clone $stashRepoUrl
		}
		else
		{
			write-host 'Have not taken care of http protocol'
		}
	}
}
popd


<#
Pretty version of http://stash/rest/api/1.0/projects/
Have removed all but the first repo - too much noise
{
    "size":4,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "key":"API",
            "id":104,
            "name":"ApiTeam",
            "public":false,
            "type":"NORMAL",
            "link":{
                "url":"/projects/API",
                "rel":"self"
            },
            "links":{
                "self":[
                    {
                        "href":"http://stash/projects/API"
                    }
                ]
            }
        },
		// More projects .... Removed
    ],
    "start":0
}


Pretty version of http://stash/rest/api/1.0/projects/ser/repos
Have removed all but the first repo - too much noise
{
    "size":17,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "slug":"travelrepublic.adverts.service",
            "id":324,
            "name":"TravelRepublic.Adverts.Service",
            "scmId":"git",
            "state":"AVAILABLE",
            "statusMessage":"Available",
            "forkable":true,
            "project":{
                "key":"SER",
                "id":121,
                "name":"Services",
                "description":"Travel Republic Services",
                "public":false,
                "type":"NORMAL",
                "link":{
                    "url":"/projects/SER",
                    "rel":"self"
                },
                "links":{
                    "self":[
                        {
                            "href":"http://stash/projects/SER"
                        }
                    ]
                }
            },
            "public":false,
            "link":{
                "url":"/projects/SER/repos/travelrepublic.adverts.service/browse",
                "rel":"self"
            },
            "cloneUrl":"http://pmcgrath@stash/scm/ser/travelrepublic.adverts.service.git",
            "links":{
                "clone":[
                    {
                        "href":"http://pmcgrath@stash/scm/ser/travelrepublic.adverts.service.git",
                        "name":"http"
                    },
                    {
                        "href":"ssh://git@stash:7999/ser/travelrepublic.adverts.service.git",
                        "name":"ssh"
                    }
                ],
                "self":[
                    {
                        "href":"http://stash/projects/SER/repos/travelrepublic.adverts.service/browse"
                    }
                ]
            }
        },
		// More repos .... Removed
	]
    "start":0
}
#>
*/
