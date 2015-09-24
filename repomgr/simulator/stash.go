/*
https://www.youtube.com/watch?v=wlR5gYd6um0
	vimtutor 	run from command line

	ctrl-e		scroll down
	ctrl-y 		scroll up
	ctrl-f		scrowl down 1 page
	ctrl-b		scrowl up 1 page
	gg		go to top of the file
	G		go to bottom of file

	w		move to next word
	e 		move to end of word

	diw		delete current word
	ciw		change current word
	cit		change current tag

	caw		change all word - includes spaces ?
	da(		delete all includes the ( and )
	dt_		delete till next space
	df_		dlete up to and including the space

	p		paste in line below
	P		paste in line above

	di"		delete inner quote, inner means current
	dip		delete inner paragraph
	. 		repeat previous command

	fa		find next "a" in line
	cta		change up to next a - deletes and puts you in insert mode
	c/ab		Change up to the next "ab" - change - find next "ab", delete text and put in insert mode

	Learning Vim as a lnguage http://benmccormick.org/2014/07/02/learning-vim-in-2014-vim-as-language/
	Vim text objects	http://blog.carbonfive.com/2011/10/17/vim-text-objects-the-definitive-guide/
https://medium.com/@mkozlows/why-atom-cant-replace-vim-433852f4b4d1
http://usevim.com/
	Your problem with vim is you don't grok vi	https://gist.github.com/nifl/1178878
	vim relative number - line numbers


	plugins based on https://www.youtube.com/watch?v=BhwtnCaFTFk
		https://github.com/bronson/vim-trailing-whitespace/tree/master/plugin
		https://github.com/scrooloose/syntastic
		https://github.com/ervandew/supertab
		https://github.com/tpope/vim-fugitive   - git in vim
		https://github.com/kien/ctrlp.vim


		https://github.com/tmuxinator/tmuxinator
		https://github.com/square/maximum-awesome


	See https://developer.atlassian.com/static/rest/stash/3.11.2/stash-rest.html
*/
package main

import (
	"log"
	"net/http"
)

var cannedResponses map[string][]byte

func main() {
	setCannedResponses()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("About to serve : %s %s\n", r.Method, r.URL.Path)
		status := 200
		defer func() { log.Printf("Served : %s %s %d\n", r.Method, r.URL.Path, status) }()

		cannedResponseKey := ""
		if r.URL.Path == "/rest/api/1.0/projects" {
			cannedResponseKey = "projects"
		} else if r.URL.Path == "/rest/api/1.0/projects/SER/repos" {
			cannedResponseKey = "projectRepositories"
		}

		if cannedResponseKey == "" {
			status = 404
			http.Error(w, http.StatusText(404), 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(cannedResponses[cannedResponseKey]))
		if err != nil {
			status = 500
			http.Error(w, err.Error(), 500)
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setCannedResponses() {
	cannedResponses = make(map[string][]byte)

	cannedResponses["projects"] = []byte(`
{
    "size":4,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "key":"API",
            "id":104,
            "name":"ApiTeam",
            "public":false,
            "type":"NORMAL",
            "link":{
                "url":"/projects/API",
                "rel":"self"
            },
            "links":{
                "self":[
                    {
                        "href":"http://stash/projects/API"
                    }
                ]
            }
        }
    ],
    "start":0
}
`)

	cannedResponses["projectRepositories"] = []byte(`
{
    "size":17,
    "limit":25,
    "isLastPage":true,
    "values":[
        {
            "slug":"travelrepublic.adverts.service",
            "id":324,
            "name":"TravelRepublic.Adverts.Service",
            "scmId":"git",
            "state":"AVAILABLE",
            "statusMessage":"Available",
            "forkable":true,
            "project":{
                "key":"SER",
                "id":121,
                "name":"Services",
                "description":"Travel Republic Services",
                "public":false,
                "type":"NORMAL",
                "link":{
                    "url":"/projects/SER",
                    "rel":"self"
                },
                "links":{
                    "self":[
                        {
                            "href":"http://stash/projects/SER"
                        }
                    ]
                }
            },
            "public":false,
            "link":{
                "url":"/projects/SER/repos/travelrepublic.adverts.service/browse",
                "rel":"self"
            },
            "cloneUrl":"http://pmcgrath@stash/scm/ser/travelrepublic.adverts.service.git",
            "links":{
                "clone":[
                    {
                        "href":"http://pmcgrath@stash/scm/ser/travelrepublic.adverts.service.git",
                        "name":"http"
                    },
                    {
                        "href":"ssh://git@stash:7999/ser/travelrepublic.adverts.service.git",
                        "name":"ssh"
                    }
                ],
                "self":[
                    {
                        "href":"http://stash/projects/SER/repos/travelrepublic.adverts.service/browse"
                    }
                ]
            }
        }
    ],
    "start":0
}
`)
}
