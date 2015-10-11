package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	commandFns := getCommandFns()

	var cmd command
	if len(os.Args) > 1 {
		cmd = command(strings.ToLower(os.Args[1]))
	}

	commandFn, ok := commandFns[cmd]
	if !ok {
		log.Fatalf("Unknown command [%s]\n", cmd)
	}

	if err := commandFn(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}
