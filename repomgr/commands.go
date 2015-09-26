/*
	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html
*/
package main

import (
	"flag"
	"log"
	"os"
)

const (
	envVarKeyPassword  = "REPO_PASSWORD"
	envVarKeyUrlPrefix = "REPO_HOST_URL"
)

func getCommandMap() map[string]func([]string) error {
	return map[string]func([]string) error{
		"github-list":       githubList,
		"stash-list":        stashList,
		"stash-multi-clone": stashMultiClone,
	}
}

func githubList(args []string) error {
	log.Println("About to run [github-list] command")

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", "https://api.github.com", "Github api url - prefix")
	userName := cmdFlags.String("username", currentUserName, "Github username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarKeyPassword), "Github password - if not supplied will be try to use an environment variable")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	log.Printf("\turl = %s\n", *url)
	log.Printf("\tuserName = %s\n", *userName)
	log.Printf("\tpassword = %s\n", *password)

	connAttrs := connectionAttributes{
		Url:        *url,
		Username:   *userName,
		Password:   *password,
		ParentName: *userName,
	}

	repoDetails, err := getGithubRepoDetails(connAttrs)
	if err != nil {
		return err
	}
	for _, repo := range repoDetails {
		log.Printf("-->%#v\n", repo)
	}

	return nil
}

func stashList(args []string) error {
	log.Println("About to run [stash-list] command")

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarKeyUrlPrefix), "Stash url - prefix")
	userName := cmdFlags.String("username", currentUserName, "Stash username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarKeyPassword), "Stash password - if not supplied will be try to use an environment variable")
	projectKey := cmdFlags.String("projectkey", "", "Stash project key")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	log.Printf("\turl = %s\n", *url)
	log.Printf("\tuserName = %s\n", *userName)
	log.Printf("\tpassword = %s\n", *password)
	log.Printf("\tprojectkey = %s\n", *projectKey)

	connAttrs := connectionAttributes{
		Url:        *url,
		Username:   *userName,
		Password:   *password,
		ParentName: *projectKey,
	}

	repoDetails, err := getStashRepoDetails(connAttrs)
	if err != nil {
		return err
	}
	for _, repo := range repoDetails {
		log.Printf("-->%#v\n", repo)
	}

	return nil
}

func stashMultiClone(args []string) error {
	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	//	url := cmdFlags.String("url", os.Getenv(envVarKeyUrlPrefix), "Stash url - prefix")
	//	userName := cmdFlags.String("username", currentUserName, "Stash username - if not supplied will be the current user's name")
	//	password := cmdFlags.String("password", os.Getenv(envVarKeyPassword), "Stash password - if not supplied will be try to use an environment variable")
	//	sshUrl := cmdFlags.String("sshurl", "ssh://git@stash:7999", "Stash ssh url - prefix")
	//	projectKey := cmdFlags.String("projectkey", "SER", "Stash project key")
	//	projectsRootDirectoryPath := flag.String("projectsrootdirectorypath", "c:/repos/stash", "Projects root directory path")
	loud := cmdFlags.Bool("loud", false, "loud flag")
	name := cmdFlags.String("name", "", "name")

	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	log.Println("About to run [stash-multi-clone] command")
	log.Printf("\tloud = %t\n", *loud)
	log.Printf("\tname = %s\n", *name)

	log.Println(currentUserName)

	return nil
}
