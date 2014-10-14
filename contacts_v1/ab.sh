#!/usr/bin/env bash

app_address=http://localhost:8080
user_name=pmcgrath
password=pass
cookie_jar=/tmp/ab-cookie-jar

# Make root call - this will give us a session
curl $app_address --cookie-jar $cookie_jar >> /dev/null
session_id=$(grep SessionId $cookie_jar | cut -f 7)
session_id_cookie="SessionId=$session_id"

# Login
curl $app_address/api/v1/login --cookie $session_id_cookie --cookie-jar $cookie_jar -XPOST -d "{ \"UserName\": \"$user_name\", \"Password\": \"$password\" }" -vvv

# See http://httpd.apache.org/docs/current/programs/ab.html
# sudo install -y apache2-utils # Only need utils do not need apache server
ab -n 100 -c 10 -C $session_id_cookie $app_address/api/v1/contacts/$user_name
