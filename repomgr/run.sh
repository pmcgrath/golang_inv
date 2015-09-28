#!/usr/bin/env bash
set -e

go build

./repomgr list -verbose -provider github -parentName pmcgrath -url https://api.github.com
./repomgr list -verbose -provider stash -parentName SER -url http://localhost:8080

./repomgr clone -verbose -provider github -parentName pmcgrath -url https://api.github.com -usessh -projectsdirectorypath /tmp/repos

