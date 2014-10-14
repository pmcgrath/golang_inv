#!/usr/bin/env bash

# Defaults - default is not to use redis, to use redis you must use the -r flag
use_redis=false
redis_address=:6379
redis_password=
session_timeout_in_minutes=1

# See http://wiki.bash-hackers.org/howto/getopts_tutorial and 
while getopts ra:p:s: opt; do
  case $opt in
    r) use_redis=true ;;
    a) redis_address=$OPTARG ;;
    p) redis_password=$OPTARG ;;
    s) session_timeout_in_minutes=$OPTARG ;;
  esac
done

# Use local redis
if $use_redis; then
  echo "Exporting for redis usage address is [$redis_address]"
  export REDIS_ADDRESS=$redis_address
  export REDIS_PASSWORD=$redis_Password

  # Ensure redis is running locally
  if [ "$(pidof redis-server)" == "" ]; then 
    echo Starting redis in background 
    [ ! -d /tmp/contacts ] && mkdir /tmp/contacts
    redis-server --dir /tmp/contacts --dbfilename contacts.rdp &
    
    # Ensure we have one user's data - no contacts
    sleep 2
    redis-cli hmset user:pmcgrath FirstName Pat LastName McGrath Password pass ContactsAsJson '[{"Id": "pmcgrath", "FirstName": "Ted", "LastName": "Toe"}]'
  fi
fi

# Session timeout
export WEBAPP_SESSION_TIMEOUT_IN_MINUTES=$session_timeout_in_minutes

# Run in background
# Can't use go run app.go as multiple files needed
# Rather than listing each time, i get all files excluding the test files passing to go run
go run $(ls *.go | grep -v _test.go) &
app_ppid=$!
echo "** app_ppid is $app_ppid - to see pstree run this : pstree -pas $app_ppid"

# Watch for file changes excluding some entries (.git directory and temp vim files) and re-run app on changes
# See http://www.alexedwards.net/blog/golang-automatic-reloads
#     http://stackoverflow.com/questions/10300835/too-many-inotify-events-while-editing-in-vim
#     http://stackoverflow.com/questions/10527936/using-inotify-to-keep-track-of-all-files-in-a-system
inotifywait -q -m -r -e close_write --exclude '(.git|.*.sw[pwx]|4913)' . | while read change_info; do 
  echo -e "\n\n** Change detected [$change_info] - re-running app" 
  notify-send "App change detected !" "\tStopping process tree $app_ppid\n\tRestarting app instance"

  pkill -TERM -P $app_ppid
  go run $(ls *.go | grep -v _test.go) &
  app_ppid=$!
  echo "** app_ppid is $app_ppid - to see pstree run this : pstree -pas $app_ppid"
done
