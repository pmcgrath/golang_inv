#!/usr/bin/env bash
set -e

go build

#./repomgr list -verbose -provider github -parentname pmcgrath -url https://api.github.com
#./repomgr list -verbose -provider stash -parentname SER -url http://localhost:8080

#./repomgr clone -verbose -provider github -parentname pmcgrath -url https://api.github.com -usessh -projectsdirectorypath /tmp/repos


#./repomgr status -verbose -projectsdirectorypath ~/oss/github.com/pmcgrath

./repomgr clone -verbose -provider github -parentname bstack -url https://api.github.com -usessh -projectsdirectorypath /tmp/repos
#./repomgr fetch -verbose -projectsdirectorypath /tmp/repos
#./repomgr remote -verbose -projectsdirectorypath /tmp/repos
#./repomgr pull -verbose -projectsdirectorypath /tmp/repos
#./repomgr status -verbose -projectsdirectorypath /tmp/repos
#./repomgr branch -verbose -projectsdirectorypath /tmp/repos
#./repomgr fetch -verbose -projectsdirectorypath /tmp/repos -remotename ted

for action in remote branch status; do
	echo -e "\n\n$action"
	./repomgr $action -projectsdirectorypath /tmp/repos
done

for action in fetch pull; do
	echo -e "\n\n$action"
	./repomgr $action -projectsdirectorypath /tmp/repos -remotename upstream
done
