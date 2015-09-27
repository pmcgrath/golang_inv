#!/usr/bin/env bash
set -e

go build

./repomgr list -provider github -parentName pmcgrath -url https://api.github.com
./repomgr list -provider stash -parentName SER -url http://localhost:8080
