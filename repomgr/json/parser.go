package main

import (
	"bytes"
	"encoding/json"
	"log"
)

var cannedResponses map[string][]byte

func main() {
	setCannedResponses()

	var data struct {
		Size       int
		IsLastPage bool
		Values     []map[string]interface{}
	}
	rdr := bytes.NewBuffer(cannedResponses["projects"])
	if err := json.NewDecoder(rdr).Decode(&data); err != nil {
		log.Fatal(err)
	}

	log.Printf("****data : %t\n", data.IsLastPage)
	for _, value := range data.Values {
		log.Printf("****Value : %#v\n", value)
		processDictionary(value)
	}
}

func processDictionary(dict map[string]interface{}) {
	log.Printf("\t-->%s\n", dict["name"].(string))
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
