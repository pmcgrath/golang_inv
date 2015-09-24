// Command line parsing
// See http://stackoverflow.com/questions/24504024/defining-independent-flagsets-in-golang
//
// Sample usages
//	.\cli-sample.go apply --silent
//
package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	envVarKeyPassword  = "REPO_PASSWORD"
	envVarKeyUrlPrefix = "REPO_URL_PREFIX"
)

func main() {
	commands := map[string]func([]string) error{
		"stash-list": stashList,
		"stash-mget": stashMGet,
	}

	command := ""
	if len(os.Args) > 1 {
		command = strings.ToLower(os.Args[1])
	}

	commandFunc, ok := commands[command]
	if !ok {
		log.Printf("Unknown command [%s]\n", command)
		return
	}

	if err := commandFunc(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

func stashList(args []string) error {
	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarKeyUrlPrefix), "Stash url - prefix")
	userName := cmdFlags.String("username", currentUserName, "Stash username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarKeyPassword), "Stash password - if not supplied will be try to use an environment variable")
	sshUrl := cmdFlags.String("sshurl", "ssh://git@stash:7999", "Stash ssh url - prefix")
	projectKey := cmdFlags.String("projectkey", "SER", "Stash project key")
	projectsRootDirectoryPath := flag.String("projectsrootdirectorypath", "c:/repos/stash", "Projects root directory path")

	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	log.Println("About to run [stash-list] command")
	log.Printf("\turl = %s\n", *url)
	log.Printf("\tuserName = %s\n", *userName)
	log.Printf("\tpassword = %s\n", *password)
	log.Printf("\tsshUrl = %s\n", *sshUrl)
	log.Printf("\tprojectkey = %s\n", *projectKey)
	log.Printf("\tprojectsRootDirectoryPath = %s\n", *projectsRootDirectoryPath)

	return nil
}

func stashMGet(args []string) error {
	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	loud := cmdFlags.Bool("loud", false, "loud flag")
	name := cmdFlags.String("name", "", "name")

	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	log.Println("About to run [stash-mget] command")
	log.Printf("\tloud = %t\n", *loud)
	log.Printf("\tname = %s\n", *name)

	return nil
}

func getCurrentUserName() (name string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	name = strings.ToLower(user.Username)
	domainSeperatorIndex := strings.Index(name, "\\")
	if domainSeperatorIndex > -1 {
		name = name[domainSeperatorIndex+1:]
	}

	return
}
