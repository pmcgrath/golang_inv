package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type TestHandler struct {
}

func (h *TestHandler) Get(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	fmt.Fprintln(w, "Success")
}

func TestRouterAddSuccess(t *testing.T) {
	spec := &Spec{t}

	handler := &TestHandler{}
	router := NewRouter()
	err := router.Add(`^/p1/?$`, handler)

	spec.Assert(err == nil, "Unexpected error : %s", err)
	spec.Assert(len(router.config) == 1, "Not added")
}

func TestRouterAddMultipleSuccess(t *testing.T) {
	spec := &Spec{t}

	handler := &TestHandler{}
	router := NewRouter()

	err := router.Add(`^/p1/?$`, handler)
	spec.Assert(err == nil, "Unexpected error : %s", err)
	spec.Assert(len(router.config) == 1, "Not added")

	err = router.Add(`^/p2/\w+/?$`, handler)
	spec.Assert(err == nil, "Unexpected error : %s", err)
	spec.Assert(len(router.config) == 2, "Not added")
}

func TestRouterAddFailureNoMethods(t *testing.T) {
	spec := &Spec{t}

	handler := 1
	router := NewRouter()
	err := router.Add(`^/p1/?$`, handler)

	spec.Assert(err != nil, "Expected error")
}

func TestRouterServeHTTPMatchFound(t *testing.T) {
	spec := &Spec{t}

	handler := &TestHandler{}
	router := NewRouter()
	router.Add(`^/p1/?$`, handler)

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/p1/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)
	spec.Assert(strings.Contains(response.Body.String(), "Success"), "Response body did not contain expected content")
}

func TestRouterServeHTTPPathMatchNotFound(t *testing.T) {
	spec := &Spec{t}

	handler := &TestHandler{}
	router := NewRouter()
	router.Add(`^/p1/?$`, handler)

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/p1/aaa", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request, requestContext)

	spec.Assert(response.Code == http.StatusNotFound, "Unexpected status code %d", response.Code)
}

func TestRouterServeHTTPMethodNotSupported(t *testing.T) {
	spec := &Spec{t}

	handler := &TestHandler{}
	router := NewRouter()
	router.Add(`^/p1/?$`, handler)

	requestContext := GetLoggedInRequestContext()
	postData := []byte("{\"Name\": \"Ted\"}")
	request, _ := http.NewRequest("POST", "/p1/", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request, requestContext)

	spec.Assert(response.Code == http.StatusMethodNotAllowed, "Unexpected status code %d", response.Code)
}
