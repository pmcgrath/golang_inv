package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

const (
	envVarNamePassword = "REPO_PASSWORD"
	envVarNameHostUrl  = "REPO_HOST_URL"
)

type commandFn func([]string) error

func getCommandFns() map[string]commandFn {
	return map[string]commandFn{
		"list": list,
		"mget": mGet,
	}
}

func list(args []string) error {
	log.Println("About to run [list] command")

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarNameHostUrl), "Host url - prefix - if not supplied will be try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentName", "", "Parent name - github organisation\\user, stash project key")
	format := cmdFlags.String("format", `{{printf "%-25s%-60s " .ParentName .Name}}{{range $key, $value := .ProtocolUrls}}{{$key}}: {{$value}} {{end}}`, "Format string for outputing the list")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	connAttrs := connectionAttributes{
		Url:      *url,
		Username: *userName,
		Password: *password,
	}

	provider, err := newProvider(*providerName, connAttrs)
	if err != nil {
		return err
	}

	repos, err := provider.getRepos(*parentName)
	if err != nil {
		return err
	}

	templateText := fmt.Sprintf(`{{range .}}%s{{println}}{{end}}`, *format)
	tmpl, err := template.New("report").Parse(templateText)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, repos)
	if err != nil {
		return err
	}

	return nil
}

func mGet(args []string) error {
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
