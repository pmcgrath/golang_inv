package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	envVarNamePassword = "REPO_PASSWORD"
	envVarNameHostUrl  = "REPO_HOST_URL"
)

type commandFn func([]string) error

type runGitCommandOnExistingReposFn func([]string) gitCmdResults

func getCommandFns() map[string]commandFn {
	return map[string]commandFn{
		"branch": branch,
		"clone":  clone,
		"fetch":  fetch,
		"list":   list,
		"pull":   pull,
		"status": status,
	}
}

func branch(args []string) error {
	log.Println("About to run [batch] command")

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(*projectsDirectoryPath,
		func(repoPaths []string) gitCmdResults {
			return execGitBranch(repoPaths)
		})
}

func clone(args []string) error {
	log.Println("About to run [clone] command")

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarNameHostUrl), "Host url - prefix - if not supplied will try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentname", "", "Parent name - github organisation\\user, stash project key")
	useSsh := cmdFlags.Bool("usessh", true, "Clone using ssh")
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	remoteName := cmdFlags.String("remotename", "upstream", "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	repos, err := getProviderRepos(*providerName, *url, *userName, *password, *parentName)
	if err != nil {
		return err
	}

	logDebugf("About to determine repos to clone, candidate count is %d\n", len(repos))
	var repoUrls []string
	for _, repo := range repos {
		// PENDING - This needs finishing - https over http - what if neither available ?
		repoUrl := repo.ProtocolUrls["http"]
		if *useSsh {
			repoUrl = repo.ProtocolUrls["ssh"]
		} else {
			// What about https
			// Include password so we do not get prompted
			// SECURITY !!!! - Over http and in logging messages
			replace := fmt.Sprintf("http://%s@", *userName)
			replaceWith := fmt.Sprintf("http://%s:%s@", *userName, *password)
			repoUrl = strings.Replace(repoUrl, replace, replaceWith, -1)
		}

		repoDirectoryName := getGitRepoNameFromUrl(repoUrl)
		repoDirectoryPath := path.Join(*projectsDirectoryPath, repoDirectoryName)
		repoAlreadyExists := testIfDirectoryExists(repoDirectoryPath)

		if !repoAlreadyExists {
			repoUrls = append(repoUrls, repoUrl)
		}
	}

	logDebugf("About to start cloning repos, count is %d\n", len(repoUrls))
	if len(repoUrls) > 0 {
		cmdResults := execGitClone(*projectsDirectoryPath, repoUrls, *remoteName)
		displayGitCmdResults(cmdResults)
	}

	return nil
}

func fetch(args []string) error {
	log.Println("About to run [fetch] command")

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	remoteName := cmdFlags.String("remotename", "upstream", "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(*projectsDirectoryPath,
		func(repoPaths []string) gitCmdResults {
			return execGitFetch(repoPaths, *remoteName)
		})
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

	repos, err := getProviderRepos(*providerName, *url, *userName, *password, *parentName)
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

func pull(args []string) error {
	log.Println("About to run [pull] command")

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	remoteName := cmdFlags.String("remotename", "upstream", "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(*projectsDirectoryPath,
		func(repoPaths []string) gitCmdResults {
			return execGitPull(repoPaths, *remoteName)
		})

}

func status(args []string) error {
	log.Println("About to run [status] command")

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(*projectsDirectoryPath,
		func(repoPaths []string) gitCmdResults {
			return execGitStatus(repoPaths)
		})
}

func runCmdOnExistingRepos(projectsDirectoryPath string, runGitCmds runGitCommandOnExistingReposFn) error {
	candidateRepoPaths, err := getAllSubDirectoryPaths(projectsDirectoryPath)
	if err != nil {
		return err
	}
	repoPaths := filterGitReposOnly(candidateRepoPaths)

	logDebugf("About to run command on repos, count is %d, out of candidate count %d\n", len(repoPaths), len(candidateRepoPaths))
	if len(repoPaths) > 0 {
		cmdResults := runGitCmds(repoPaths)
		displayGitCmdResults(cmdResults)
	}

	return nil
}

func displayGitCmdResults(results gitCmdResults) {
	for _, result := range results {
		fmt.Printf("\x1b[32m%s\x1b[0m\n", result.RepoPath)
		if result.Error != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", result.Error)
		}
		for _, line := range result.Output {
			fmt.Printf("\t%s\n", line)
		}
		fmt.Println()
	}
}

func getProviderRepos(providerName, url, userName, password, parentName string) (repositoryDetails, error) {
	logDebugf("About to instantiate provider [%s]\n", providerName)
	connAttrs := providerConnectionAttributes{
		Url:      url,
		Username: userName,
		Password: password,
	}
	provider, err := newProvider(providerName, connAttrs)
	if err != nil {
		return nil, err
	}

	logDebugf("About to get repos for provider [%s]\n", providerName)
	repos, err := provider.getRepos(parentName)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
