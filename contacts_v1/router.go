package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type router struct {
	config map[*regexp.Regexp]*pathEntry
}

func NewRouter() *router {
	return &router{
		config: make(map[*regexp.Regexp]*pathEntry),
	}
}

func (router *router) Add(pattern string, pathHandler interface{}) error {
	// See http://stackoverflow.com/questions/20714939/how-to-properly-use-call-in-reflect-package-golang
	key := regexp.MustCompile(pattern)

	pathEntry := &pathEntry{pattern: pattern, supportedMethods: make(map[string]ContextualHandlerFunc)}

	interfaceValue := reflect.ValueOf(pathHandler)

	methodFound := false
	for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
		capitalisedMethod := strings.Title(strings.ToLower(method))

		methodValue := interfaceValue.MethodByName(capitalisedMethod)
		if !methodValue.IsValid() {
			continue
		}

		methodInterface := methodValue.Interface()
		methodHandle := methodInterface.(func(http.ResponseWriter, *http.Request, *RequestContext)) // Assuming func signature matches
		pathEntry.supportedMethods[method] = methodHandle                                           // Taking care of func(.. to HandlerFunc conversion

		methodFound = true
	}

	if methodFound != true {
		return errors.New("No method found")
	}

	router.config[key] = pathEntry
	return nil
}

func (router *router) Get(path string) *pathEntry {
	for key, value := range router.config {
		if key.MatchString(path) {
			return value
		}
	}

	return nil
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	log.Printf("%s %s %s Servicing\n", c.GetLogMessagePrefix(), r.URL.Path, r.Method)

	isPathSupported, isMethodSupported := false, false
	pathEntry := router.Get(r.URL.Path)
	if pathEntry != nil {
		isPathSupported = true
		methodHandler := pathEntry.Get(r.Method)
		if methodHandler != nil {
			isMethodSupported = true
			methodHandler(w, r, c)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	log.Printf("%s %s %s Serviced: Path supported = %t, method supported = %t\n", c.GetLogMessagePrefix(), r.URL.Path, r.Method, isPathSupported, isMethodSupported)
}

// Path entry - config
type pathEntry struct {
	pattern          string
	supportedMethods map[string]ContextualHandlerFunc
}

func (entry *pathEntry) Get(method string) ContextualHandlerFunc {
	return entry.supportedMethods[method]
}
