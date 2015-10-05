$ErrorActionPreference = 'Stop'

go build
if ($LastExitCode -ne 0) { return }

#./repomgr list -verbose -provider github -parentname pmcgrath -url https://api.github.com
#./repomgr list -verbose -provider stash -parentname SER -url http://stash

#./repomgr clone -verbose -provider stash -parentName SER -url http://stash -usessh -projectsdirectorypath c:/repos/stash/ser
./repomgr status -verbose -projectsdirectorypath c:/repos/stash/ser
#./repomgr fetch -verbose -projectsdirectorypath c:/repos/stash/ser
