package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestAssetsHandlerGet(t *testing.T) {
	spec := &Spec{t}

	handler := &AssetsHandler{}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/assets/js/main.js", nil)
	response := httptest.NewRecorder()

	handler.Get(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)
}

func TestContactApiHandlerDeleteSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("DELETE", "/api/v1/contacts/pmcgrath/pmcgrath/", nil)
	response := httptest.NewRecorder()

	handler.Delete(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)

	user, _ := store.Get("pmcgrath")
	spec.Assert(len(user.Contacts) == 1, "Unexpected contact count %d", len(user.Contacts))

}

func TestContactApiHandlerDeleteResourceDoesNotExist(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("DELETE", "/api/v1/contacts/pmcgrath/DOESNOTEXIST", nil)
	response := httptest.NewRecorder()

	handler.Delete(response, request, requestContext)

	spec.Assert(response.Code == http.StatusNotFound, "Unexpected status code %d", response.Code)

	user, _ := store.Get("pmcgrath")
	spec.Assert(len(user.Contacts) == 2, "Unexpected contact count %d", len(user.Contacts))
}

func TestContactApiHandlerGetSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/api/v1/contacts/pmcgrath/ted", nil)
	response := httptest.NewRecorder()

	handler.Get(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)

	contentTypeHeader := response.HeaderMap["Content-Type"][0]
	spec.Assert(contentTypeHeader == "application/json", "Unexpected content type header %s", contentTypeHeader)

	body := response.Body.String()
	spec.Assert(strings.Contains(body, `"Id":"ted",`), "Response body did not contain expected content, body is %s", body)
}

func TestContactApiHandlerGetNotFound(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/api/v1/contacts/pmcgrath/DOESNOTEXIST", nil)
	response := httptest.NewRecorder()

	handler.Get(response, request, requestContext)

	spec.Assert(response.Code == http.StatusNotFound, "Unexpected status code %d", response.Code)
}

func TestContactApiHandlerPutSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	postData := []byte(`{"FirstName": "Ted", "LastName": "Toad"}`)
	request, _ := http.NewRequest("PUT", "/api/v1/contacts/pmcgrath/ted", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Put(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)
}

func TestContactApiHandlerPutUrlAndBodyIdConflict(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	postData := []byte(`{"Id": "..", "FirstName": "Ted", "LastName": "Toad"}`)
	request, _ := http.NewRequest("PUT", "/api/v1/contacts/pmcgrath/ted", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Put(response, request, requestContext)

	spec.Assert(response.Code == http.StatusBadRequest, "Unexpected status code %d", response.Code)
}

func TestContactApiHandlerPutFailureDueToIncompleteData(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	postData := []byte(`{"LastName": "Toad"}`) // No first name
	request, _ := http.NewRequest("PUT", "/api/v1/contacts/pmcgrath/ted", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Put(response, request, requestContext)

	spec.Assert(response.Code == http.StatusBadRequest, "Unexpected status code %d", response.Code)
}

func TestContactsApiHandlerGetSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactsApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	request, _ := http.NewRequest("GET", "/api/v1/contacts/pmcgrath", nil)
	response := httptest.NewRecorder()

	handler.Get(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)

	body := response.Body.String()
	spec.Assert(strings.Contains(body, `"Id":"ted",`), "Response body did not contain expected content, body is %s", body)
}

func TestContactsApiHandlerPostSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactsApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	postData := []byte(`{"FirstName": "Ted", "LastName": "Toad", "Email": "tt@gmail.com"}`)
	request, _ := http.NewRequest("POST", "/api/v1/contacts/pmcgrath", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Post(response, request, requestContext)

	spec.Assert(response.Code == http.StatusCreated, "Unexpected status code %d", response.Code)

	locationHeader := response.HeaderMap["Location"][0]
	spec.Assert(strings.HasPrefix(locationHeader, request.URL.Path+"/"), "Unexpected location header %s", locationHeader)
}

func TestContactsApiHandlerPostFailureDueToInCompleteData(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &ContactsApiHandler{PathPrefix: "/api/v1/contacts/", Store: store}

	requestContext := GetLoggedInRequestContext()
	postData := []byte(`{"FirstName": "Ted", "Email": "tt@gmail.com"}`) // Missing lastname
	request, _ := http.NewRequest("POST", "/api/v1/contacts/pmcgrath", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Post(response, request, requestContext)

	spec.Assert(response.Code == http.StatusBadRequest, "Unexpected status code %d", response.Code)
}

func TestLogInApiHandlerDeleteSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &LogInApiHandler{Store: store}

	requestContext := GetLoggedInRequestContext()

	request, _ := http.NewRequest("DELETE", "/api/v1/login", nil)
	response := httptest.NewRecorder()

	handler.Delete(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)

	spec.Assert(requestContext.Session.UserName == "", "Unexpected session user name %s", requestContext.Session.UserName)
}

func TestLogInApiHandlerPostSuccess(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &LogInApiHandler{Store: store}

	requestContext := GetLoggedInRequestContext()
	requestContext.Session.UserName = ""

	postData := []byte(`{"UserName": "pmcgrath", "Password": "pass"}`)
	request, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Post(response, request, requestContext)

	spec.Assert(response.Code == http.StatusOK, "Unexpected status code %d", response.Code)

	spec.Assert(requestContext.Session.UserName == "pmcgrath", "Unexpected session user name %s", requestContext.Session.UserName)
}

func TestLogInApiHandlerPostWhereAlreadyLoggedInFailure(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &LogInApiHandler{Store: store}

	requestContext := GetLoggedInRequestContext()

	postData := []byte(`{"UserName": "pmcgrath", "Password": "pass"}`)
	request, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Post(response, request, requestContext)

	spec.Assert(response.Code == http.StatusForbidden, "Unexpected status code %d", response.Code)
}

func TestLogInApiHandlerPostIncorrectPassword(t *testing.T) {
	spec := &Spec{t}

	store := GetInitialisedUserStore()
	handler := &LogInApiHandler{Store: store}

	requestContext := GetLoggedInRequestContext()
	requestContext.Session.UserName = ""

	postData := []byte(`{"UserName": "pmcgrath", "Password": "BADPASS"}`)
	request, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewReader(postData))
	response := httptest.NewRecorder()

	handler.Post(response, request, requestContext)

	spec.Assert(response.Code == http.StatusUnauthorized, "Unexpected status code %d", response.Code)

	spec.Assert(requestContext.Session.UserName == "", "Unexpected session user name %s", requestContext.Session.UserName)
}

func GetInitialisedUserStore() UserStore {
	store := NewInMemoryUserStore()
	store.Save(
		&User{
			Id:        "pmcgrath",
			FirstName: "Pat",
			LastName:  "Mc Grath",
			Email:     "pmcgrath@gmail.com",
			Password:  "pass",
			Contacts: []Contact{
				Contact{
					Id:        "pmcgrath",
					FirstName: "Peter",
					LastName:  "Mc Grath",
					Phones: []Phone{
						Phone{
							Description: "Home",
							Number:      "44 066 7132310",
						},
					},
				},
				Contact{
					Id:        "ted",
					FirstName: "Ted",
					LastName:  "Toe",
					Phones: []Phone{
						Phone{
							Description: "Home",
							Number:      "353 066 7132310",
						},
					},
				},
			},
		})

	return store
}

func GetLoggedInRequestContext() *RequestContext {
	return &RequestContext{
		Id:        Uuid(),
		StartTime: time.Now(),
		Session: &Session{
			Id:         Uuid(),
			UserName:   "pmcgrath",
			Data:       make(map[string]interface{}),
			LastAccess: time.Now(),
		},
		Data: make(map[string]interface{}),
	}
}
