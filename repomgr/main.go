package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	//PENDING - investigate merging of streams - goroutines - while running concurrent shell outs for git
	dirs, _ := getAllSubDirectoryPaths("/home/pmcgrath/oss/github.com/pmcgrath")
	for _, dir := range dirs {
		log.Printf("D --> %s\n", dir)
	}
	repos := filterGitReposOnly(dirs)
	for _, repo := range repos {
		log.Printf("--> %s\n", repo)
	}
	return

	commandFns := getCommandFns()

	command := ""
	if len(os.Args) > 1 {
		command = strings.ToLower(os.Args[1])
	}

	commandFn, ok := commandFns[command]
	if !ok {
		log.Fatalf("Unknown command [%s]\n", command)
	}

	if err := commandFn(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}
