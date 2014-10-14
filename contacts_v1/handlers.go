package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Root handler
type RootHandler struct {
}

func (h *RootHandler) Get(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	t, _ := template.New("Html").Parse(rootHtmlTemplate)
	err := t.Execute(w, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Assets handler
type AssetsHandler struct {
}

func (h *AssetsHandler) Get(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if content, ok := assetMap[r.URL.Path]; ok {
		contentType := getAssetContentType(r.URL.Path)
		w.Header().Add("Content-Type", contentType)

		if _, err := fmt.Fprintf(w, content); err != nil {
			log.Printf("%s Error detected when trying write asset : %s\n", c.GetLogMessagePrefix(), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// Contact api handler
type ContactApiHandler struct {
	PathPrefix string
	Store      UserStore
}

func (h *ContactApiHandler) PreProcess(w http.ResponseWriter, r *http.Request, c *RequestContext) bool {
	ids, err := getIdsFromUrlPath(r.URL.Path, h.PathPrefix)
	if err != nil {
		log.Printf("%s Error detected when trying to get ids from url : %s\n", c.GetLogMessagePrefix(), err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false
	}

	userId, contactId := ids[0], ids[1]
	if userId != c.GetUserName() {
		log.Printf("%s Forbidden, context user id %s\n", c.GetLogMessagePrefix(), c.GetUserName())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}

	user, err := h.Store.Get(userId)
	if err != nil {
		log.Printf("%s Error detected when trying to get user with id %s : %s\n", c.GetLogMessagePrefix(), userId, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	c.Data["User"] = user
	c.Data["ContactId"] = contactId
	return true
}

func (h *ContactApiHandler) Delete(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !h.PreProcess(w, r, c) {
		return
	}
	user := c.Data["User"].(*User)
	contactId := c.Data["ContactId"].(string)

	index, ok := user.GetContactIndex(contactId)
	if !ok {
		log.Printf("%s Contact not found for user with id %s and contact with id %s\n", c.GetLogMessagePrefix(), user.Id, contactId)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	newContacts := append(user.Contacts[:index], user.Contacts[index+1:]...)
	user.Contacts = newContacts
	h.Store.Save(user)
}

func (h *ContactApiHandler) Get(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !h.PreProcess(w, r, c) {
		return
	}
	user := c.Data["User"].(*User)
	contactId := c.Data["ContactId"].(string)

	index, ok := user.GetContactIndex(contactId)
	if !ok {
		log.Printf("%s Contact not found for user with id %s and contact with id %s\n", c.GetLogMessagePrefix(), user.Id, contactId)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(user.Contacts[index]); err != nil {
		log.Printf("%s Error detected when trying to encode contact for user with id %s and contact with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, contactId, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ContactApiHandler) Put(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !h.PreProcess(w, r, c) {
		return
	}
	user := c.Data["User"].(*User)
	contactId := c.Data["ContactId"].(string)

	var contact Contact
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&contact)
	if err != nil {
		log.Printf("%s Error detected when trying to decode contact for user with id %s and contact with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, contactId, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if contact.Id == "" {
		contact.Id = contactId
	}
	if contact.Id != contactId {
		log.Printf("%s Contact id conflict url is %s put body is %s\n", c.GetLogMessagePrefix(), contactId, contact.Id)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if valid, err := (&contact).IsValidForSaving(); !valid {
		log.Printf("%s Contact state is not valid for saving for user with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	index, ok := user.GetContactIndex(contactId)
	if ok {
		user.Contacts[index] = contact
	} else {
		user.Contacts = append(user.Contacts, contact)
	}

	err = h.Store.Save(user)
	if err != nil {
		log.Printf("%s Error detected when saving user with id %s and contact with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, contactId, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Contacts api handler
type ContactsApiHandler struct {
	PathPrefix string
	Store      UserStore
}

func (h *ContactsApiHandler) PreProcess(w http.ResponseWriter, r *http.Request, c *RequestContext) bool {
	ids, err := getIdsFromUrlPath(r.URL.Path, h.PathPrefix)
	if err != nil {
		log.Printf("%s Error detected when trying to get ids from url : %s\n", c.GetLogMessagePrefix(), err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return false
	}

	userId := ids[0]
	if userId != c.GetUserName() {
		log.Printf("%s Forbidden, context user id %s\n", c.GetLogMessagePrefix(), c.GetUserName())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}

	user, err := h.Store.Get(userId)
	if err != nil {
		log.Printf("%s Error detected when trying to get user with id %s : %s\n", c.GetLogMessagePrefix(), userId, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	c.Data["User"] = user
	return true
}

func (h *ContactsApiHandler) GenerateUrl(userId, contactId string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(h.PathPrefix, "/"), userId, contactId)
}

func (h *ContactsApiHandler) Get(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !h.PreProcess(w, r, c) {
		return
	}
	user := c.Data["User"].(*User)

	contacts := make([]Contact, 0)
	if contacts != nil {
		contacts = user.Contacts
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(contacts); err != nil {
		log.Printf("%s Error detected when trying to encode contacts for user with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *ContactsApiHandler) Post(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if !h.PreProcess(w, r, c) {
		return
	}
	user := c.Data["User"].(*User)

	var contact Contact
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&contact)
	if err != nil {
		log.Printf("%s Error detected when trying to decode contact for user with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	contact.Id = Uuid()
	if valid, err := (&contact).IsValidForSaving(); !valid {
		log.Printf("%s Contact state is not valid for saving for user with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user.Contacts = append(user.Contacts, contact)

	err = h.Store.Save(user)
	if err != nil {
		log.Printf("%s Error detected when saving user with id %s : %s\n", c.GetLogMessagePrefix(), user.Id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contactUrl := h.GenerateUrl(user.Id, contact.Id)

	w.Header().Set("Location", contactUrl)
	w.WriteHeader(http.StatusCreated)
}

// LogIn api handler
type LogInApiHandler struct {
	Store UserStore
}

func (h *LogInApiHandler) Delete(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if c.Session.UserName == "" {
		log.Printf("%s User not logged in\n", c.GetLogMessagePrefix())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	c.Session.UserName = ""
}

func (h *LogInApiHandler) Post(w http.ResponseWriter, r *http.Request, c *RequestContext) {
	if c.Session.UserName != "" {
		log.Printf("%s User %s already logged in, must log out first\n", c.GetLogMessagePrefix(), c.Session.UserName)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var credentials map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)
	if err != nil {
		log.Printf("%s Error detected when trying to decode credentials : %s\n", c.GetLogMessagePrefix(), err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId, ok := credentials["UserName"]
	if !ok {
		log.Printf("%s User name not suppplied\n", c.GetLogMessagePrefix())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	password, ok := credentials["Password"]
	if !ok {
		log.Printf("%s Password not suppplied\n", c.GetLogMessagePrefix())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := h.Store.Get(userId)
	if err != nil {
		log.Printf("%s User record not found for user id %s\n", c.GetLogMessagePrefix(), userId)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !user.Authenticate(password) {
		log.Printf("%s User password is incorrect for user id %s\n", c.GetLogMessagePrefix(), userId)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	c.Session.UserName = user.Id
}
