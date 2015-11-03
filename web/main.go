package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get list of services
	services := []struct {
		Name string
	}{
		{Name: "Svc1"},
		{Name: "Svc2"},
	}

	err := indexTemplate.Execute(w, services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiServices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get list of services and render as json
	services := []struct {
		Name string
	}{
		{Name: "Svc1"},
		{Name: "Svc2"},
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(services)
}

func documentation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "documentation\n")
}

func resources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "resources for %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/api/services", apiServices)
	router.GET("/documentation", documentation)
	router.GET("/resources", resources)
	router.GET("/resources/:name", resources)

	log.Fatal(http.ListenAndServe(":8080", router))
}
