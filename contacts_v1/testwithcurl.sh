block=true

function reportstatus() {
	echo -e ">>>>> $1"
	if $block; then read; fi 
}

cookie_file_path=/tmp/contact_cookies.txt

curl http://localhost:8080/api/v1/login --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv -XPOST -d '{ "UserName": "pmcgrath", "Password": "pass" }'
reportstatus "Expected a 200 - logged in"

curl http://localhost:8080/api/v1/contacts/pmcgrath --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv -XPOST -d '{ "FirstName": "Tom", "LastName": "Toe" }'
reportstatus "Expected a 201 - contact created"

curl http://localhost:8080/api/v1/contacts/pmcgrath/teddydoesnotexist --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv -XDELETE
reportstatus "Expected a 404 - contact does not exist so cannot be deleted"

curl http://localhost:8080/api/v1/contacts/pmcgrath --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 200 - contacts found"

curl http://localhost:8080/api/v1/contacts/tedtoecontacts --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 403 - cannot get someone else's contacts"

curl http://localhost:8080/tedtoepath --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 404 - unknown resource - url is unexpected"

curl http://localhost:8080/api/v1/contacts/pmcgrath --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 200 - contacts found"

rm $cookie_file_path
curl http://localhost:8080/assets/ted --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 200 - asset resource found, do not need to be logged in"

curl http://localhost:8080/assets/ted --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv -XPOST
reportstatus "Expected a 405 - cannot post to an asset resource"

curl http://localhost:8080/api/v1/contacts/pmcgrath --cookie $cookie_file_path --cookie-jar $cookie_file_path -vv
reportstatus "Expected a 401 - user is not logged in - we deleted the cookie so we are using a session where no user logged in"
