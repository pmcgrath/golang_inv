package main

import (
	"bytes"
	"encoding/json"
	"log"
	"runtime"
)

var cannedResponses map[string][]byte

func init() {
	setCannedResponses()
}

func main() {
	parseProjectsWithExplicitType()
	parseProjectsWithGenericType()

	parseProjectRepoitoriesWithExplicitType()
	parseProjectRepoitoriesWithGenericType()

	parseBothWithiCombinedType()
}

func parseProjectsWithExplicitType() {
	// Examined the json and used golang data structures to match the json data
	var data struct {
		Size       int
		IsLastPage bool
		Values     []struct {
			Name string
			Link struct {
				Url string
				Rel string
			}
		}
	}
	decodeJson("projects", &data)

	log.Println(getFuncName())
	for _, value := range data.Values {
		log.Printf("\tName: %s\tUrl: %s\n", value.Name, value.Link.Url)
	}
}

func parseProjectsWithGenericType() {
	// Here we defined a type for the parent object based on the json, we then use generic structures to get to the data
	// Needed to use printf with "%#v" to decide what data structures to use
	var data struct {
		Size       int
		IsLastPage bool
		Values     []map[string]interface{}
	}
	decodeJson("projects", &data)

	log.Printf("\n\n%s\n", getFuncName())
	for _, value := range data.Values {
		name := value["name"].(string)
		link := value["link"].(map[string]interface{})
		url := link["url"].(string)

		log.Printf("\tName: %s\tUrl: %s\n", name, url)
	}
}

func parseProjectRepoitoriesWithExplicitType() {
	// Examined the json and used golang data structures to match the json data
	var data struct {
		Size       int
		IsLastPage bool
		Values     []struct {
			Name  string
			Links struct {
				Clone []struct {
					Href string
					Name string
				}
			}
		}
	}
	decodeJson("projectRepositories", &data)

	log.Printf("\n\n%s\n", getFuncName())
	for _, value := range data.Values {
		log.Printf("\tName: %s\n", value.Name)
		for _, clone := range value.Links.Clone {
			log.Printf("\t\tName: %s Href: %s\n", clone.Name, clone.Href)
		}
	}
}

func parseProjectRepoitoriesWithGenericType() {
	// Here we defined a type for the parent object based on the json, we then use generic structures to get to the data
	// Needed to use printf with "%#v" to decide what data structures to use
	var data struct {
		Size       int
		IsLastPage bool
		Values     []map[string]interface{}
	}
	decodeJson("projectRepositories", &data)

	log.Printf("\n\n%s\n", getFuncName())
	for _, value := range data.Values {
		name := value["name"].(string)
		links := value["links"].(map[string]interface{})
		clones := links["clone"].([]interface{})

		log.Printf("\tName: %s\n", name)
		for _, clone := range clones {
			cloneMap := clone.(map[string]interface{})
			name := cloneMap["name"].(string)
			href := cloneMap["href"].(string)

			log.Printf("\t\tName: %s Href: %s\n", name, href)
		}
	}
}

func parseBothWithiCombinedType() {
	// Examined the json and used golang data structures to match the json data, since both have no conflicts we could union here
	// The decode does not zero out the data so you can get data from a previous call, if no matching data in the subsequent json will remain in the dtaa structure
	var data struct {
		Size       int
		IsLastPage bool
		Values     []struct {
			Name string
			Link struct {
				Url string
				Rel string
			}
			Links struct {
				Clone []struct {
					Href string
					Name string
				}
			}
		}
	}

	log.Printf("\n\n%s\n", getFuncName())

	decodeJson("projects", &data)
	for _, value := range data.Values {
		log.Printf("\tName: %s\tUrl: %s\n", value.Name, value.Link.Url)
	}

	decodeJson("projectRepositories", &data)
	for _, value := range data.Values {
		log.Printf("\tName: %s\n", value.Name)
		for _, clone := range value.Links.Clone {
			log.Printf("\t\tName: %s Href: %s\n", clone.Name, clone.Href)
		}
	}
}

func decodeJson(responseKey string, respData interface{}) {
	buf := bytes.NewBuffer(cannedResponses[responseKey])
	if err := json.NewDecoder(buf).Decode(respData); err != nil {
		log.Fatal(err)
	}
}

func getFuncName() string {
	// See 	http://stackoverflow.com/questions/10742749/get-name-of-function-using-google-gos-reflection
	//	http://play.golang.org/p/teu5CnHoek
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	return runtime.FuncForPC(pc).Name()
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
