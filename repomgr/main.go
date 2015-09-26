package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	commands := getCommandMap()

	command := ""
	if len(os.Args) > 1 {
		command = strings.ToLower(os.Args[1])
	}

	commandFunc, ok := commands[command]
	if !ok {
		log.Fatalf("Unknown command [%s]\n", command)
	}

	if err := commandFunc(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}
