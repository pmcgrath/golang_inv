package main

import (
	"log"
	"os"
	"strings"
)

func main() {
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
