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
	url := cmdFlags.String("url", os.Getenv(envVarNameHostUrl), "Host url - prefix - if not supplied will try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentName", "", "Parent name - github organisation\\user, stash project key")
	format := cmdFlags.String("format", `{{printf "%-25s%-60s " .ParentName .Name}}{{range $key, $value := .ProtocolUrls}}{{$key}}: {{$value}} {{end}}`, "Format string for outputing the list")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	logDebugf("About to instantiate provider [%s]\n", *providerName)
	connAttrs := connectionAttributes{
		Url:      *url,
		Username: *userName,
		Password: *password,
	}
	provider, err := newProvider(*providerName, connAttrs)
	if err != nil {
		return err
	}

	logDebugf("About to get repos for provider [%s]\n", *providerName)
	repos, err := provider.getRepos(*parentName)
	if err != nil {
		return err
	}

	logDebugf("About to list repos for provider [%s] - %d repos found\n", *providerName, len(repos))
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
	log.Println("About to run [mget] command")

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarNameHostUrl), "Host url - prefix - if not supplied will try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentName", "", "Parent name - github organisation\\user, stash project key")
	useSsh := cmdFlags.Bool("usessh", true, "Clone using ssh")
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	logDebugf("About to instantiate provider [%s]\n", *providerName)
	connAttrs := connectionAttributes{
		Url:      *url,
		Username: *userName,
		Password: *password,
	}
	provider, err := newProvider(*providerName, connAttrs)
	if err != nil {
		return err
	}

	logDebugf("About to get repos for provider [%s]\n", *providerName)
	repos, err := provider.getRepos(*parentName)
	if err != nil {
		return err
	}

	logDebugf("About to start cloning repos, count is %d\n", len(repos))
	os.Chdir(*projectsDirectoryPath)
	for _, repo := range repos {
		if repo.Name != "echo" || repo.Name != "ddash" {
			continue
		}

		repoUrl := repo.ProtocolUrls["http"]
		if *useSsh {
			repoUrl = repo.ProtocolUrls["ssh"]
		}

		repoPath := *projectsDirectoryPath + "/" + repo.Name

		log.Printf("git clone %s %s\n", repoUrl, repoPath)
	}

	return nil
}
