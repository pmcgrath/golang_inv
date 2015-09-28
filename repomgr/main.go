package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	//PENDING - investigate merging of streams - goroutines - while running concurrent shell outs for git
	r := filterGitReposOnly([]string{"c:/repos/stash/ser/ted", "c:/repos/stash/ser/travelrepublic.adverts.service"})
	log.Printf("Expect to see adverts svs here %#v", r)
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
