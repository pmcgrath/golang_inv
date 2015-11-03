package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

type command string

const (
	unknownCmd command = "unknown"
	branchCmd          = "branch"
	cloneCmd           = "clone"
	fetchCmd           = "fetch"
	listCmd            = "list"
	pullCmd            = "pull"
	remoteCmd          = "remote"
	statusCmd          = "status"
)

const (
	envVarNamePassword = "REPO_PASSWORD"
	envVarNameHostURL  = "REPO_HOST_URL"
)

type commandFn func([]string) error

func getCommandFns() map[command]commandFn {
	return map[command]commandFn{
		branchCmd: branch,
		cloneCmd:  clone,
		fetchCmd:  fetch,
		listCmd:   list,
		pullCmd:   pull,
		remoteCmd: remote,
		statusCmd: status,
	}
}

func branch(args []string) error {
	log.Printf("About to run [%s] command\n", branchCmd)

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(branchCmd, *projectsDirectoryPath, "")
}

func clone(args []string) error {
	log.Printf("About to run [%s] command", cloneCmd)

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarNameHostURL), "Host url - prefix - if not supplied will try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentname", "", "Parent name - github organisation\\user, stash project key")
	useSSH := cmdFlags.Bool("usessh", true, "Clone using ssh")
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
	var repoURLs []string
	for _, repo := range repos {
		repoURL := ""
		if *useSSH {
			repoURL = repo.ProtocolURLs["ssh"]
		} else {
			securityMessage := `
			// SECURITY !!!!
			// 	Not sure I want to even support this
			//	Username\password management
			//		Don't want the remote to include the password which will be stored on disk
			//		Don't want to have to support password prompts - not sure even if i could
			//		Could use a credential helper - this app could act as one and feed the password back out ???
			//			https://www.kernel.org/pub/software/scm/git/docs/git-credential-store.html
			//			http://git-scm.com/docs/git-credential-cache
			//			http://stackoverflow.com/questions/5343068/is-there-a-way-to-skip-password-typing-when-using-https-github
			//	Http - password in logs
			// 	Should be at least https
			//	Would be good to store the username as part of the remote url, should this be optional
`
			log.Println("Using http - your password is in the remote url on disk !!!\n\n", securityMessage)

			repoURL = repo.ProtocolURLs["http"]
			replace := fmt.Sprintf("http://%s@", *userName)
			replaceWith := fmt.Sprintf("http://%s:%s@", *userName, *password)
			repoURL = strings.Replace(repoURL, replace, replaceWith, -1)
		}

		repoDirectoryName := getGitRepoNameFromURL(repoURL)
		repoDirectoryPath := path.Join(*projectsDirectoryPath, repoDirectoryName)
		repoAlreadyExists := testIfDirectoryExists(repoDirectoryPath)

		if !repoAlreadyExists {
			repoURLs = append(repoURLs, repoURL)
		}
	}

	logDebugf("About to start cloning repos, count is %d\n", len(repoURLs))
	if len(repoURLs) > 0 {
		cmdResults := execGitClone(*projectsDirectoryPath, repoURLs, *remoteName)
		displayGitCmdResults(cmdResults)
	}

	return nil
}

func fetch(args []string) error {
	log.Printf("About to run [%s] command", fetchCmd)

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	remoteName := cmdFlags.String("remotename", "upstream", "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(fetchCmd, *projectsDirectoryPath, *remoteName)
}

func list(args []string) error {
	log.Printf("About to run [%s] command", listCmd)

	currentUserName, err := getCurrentUserName()
	if err != nil {
		return err
	}

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	url := cmdFlags.String("url", os.Getenv(envVarNameHostURL), "Host url - prefix - if not supplied will try to use the REPO_HOST_URL environment variable")
	providerName := cmdFlags.String("provider", "", "Provider - github, stash")
	userName := cmdFlags.String("username", currentUserName, "Username - if not supplied will be the current user's name")
	password := cmdFlags.String("password", os.Getenv(envVarNamePassword), "Password - if not supplied will be try to use the REPO_PASSWORD environment variable")
	parentName := cmdFlags.String("parentname", "", "Parent name - github organisation\\user, stash project key")
	format := cmdFlags.String("format", `{{printf "%-25s%-60s " .ParentName .Name}}{{range $key, $value := .ProtocolURLs}}{{$key}}: {{$value}} {{end}}`, "Format string for outputing the list")
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
	log.Printf("About to run [%s] command", pullCmd)

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	remoteName := cmdFlags.String("remotename", "upstream", "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(pullCmd, *projectsDirectoryPath, *remoteName)
}

func remote(args []string) error {
	log.Printf("About to run [%s] command", remoteCmd)

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(remoteCmd, *projectsDirectoryPath, "")
}

func status(args []string) error {
	log.Printf("About to run [%s] command", statusCmd)

	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	projectsDirectoryPath := cmdFlags.String("projectsdirectorypath", getDefaultProjectsDirectoryPath(), "Projects directory path")
	verbose := cmdFlags.Bool("verbose", false, "Verbose flag")
	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	isVerbose = *verbose

	return runCmdOnExistingRepos(statusCmd, *projectsDirectoryPath, "")
}

func runCmdOnExistingRepos(command command, projectsDirectoryPath, remoteName string) error {
	candidateRepoPaths, err := getAllSubDirectoryPaths(projectsDirectoryPath)
	if err != nil {
		return err
	}
	repoPaths := filterGitReposOnly(candidateRepoPaths)

	logDebugf("About to run command [%s] on repos, count is %d, out of candidate count %d\n", command, len(repoPaths), len(candidateRepoPaths))
	if len(repoPaths) > 0 {
		var cmdResults gitCmdResults
		switch command {
		case branchCmd:
			cmdResults = execGitBranch(repoPaths)
		case fetchCmd:
			cmdResults = execGitFetch(repoPaths, remoteName)
		case pullCmd:
			cmdResults = execGitPull(repoPaths, remoteName)
		case remoteCmd:
			cmdResults = execGitRemote(repoPaths)
		case statusCmd:
			cmdResults = execGitStatus(repoPaths)
		default:
			return fmt.Errorf("Unexpected command [%s]", command)
		}

		displayGitCmdResults(cmdResults)
	}

	return nil
}

func displayGitCmdResults(results gitCmdResults) {
	greenColour, redColour, resetColour := "\x1b[32m", "\x1b[31m", "\x1b[0m"
	if runtime.GOOS == "windows" {
		greenColour, redColour, resetColour = "", "", ""
	}

	for _, result := range results {
		fmt.Printf("%s%s%s\n", greenColour, result.RepoPath, resetColour)
		if result.Error != nil {
			fmt.Printf("%s%s%s\n", redColour, result.Error, resetColour)
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
		URL:      url,
		Username: userName,
		Password: password,
	}
	provider, err := newProvider(providerName, connAttrs)
	if err != nil {
		return nil, err
	}

	logDebugf("About to get repositories for provider [%s]\n", providerName)
	repos, err := provider.getRepositories(parentName)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
