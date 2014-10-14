package main

import (
	"fmt"
	"strings"
	"time"
)

type Session struct {
	Id         string
	UserName   string
	Data       map[string]interface{}
	LastAccess time.Time
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Contacts  []Contact
}

func (user *User) Authenticate(password string) bool {
	return password == user.Password // This is much too simplistic, but for now
}

func (user *User) GetContactIndex(id string) (int, bool) {
	for index, contact := range user.Contacts {
		if contact.Id == id {
			return index, true
		}
	}
	return -1, false
}

type Contact struct {
	Id        string  `json:",omitempty"`
	FirstName string  `json:",omitempty"`
	LastName  string  `json:",omitempty"`
	Emails    []Email `json:",omitempty"`
	Phones    []Phone `json:",omitempty"`
	Twitter   string  `json:",omitempty"`
	Notes     string  `json:",omitempty"`
}

type Email struct {
	Description string `json:",omitempty"`
	Address     string `json:",omitempty"`
}

type Phone struct {
	Description string `json:",omitempty"`
	Number      string `json:",omitempty"`
}

func (contact *Contact) IsValidForSaving() (bool, error) {
	err := ""
	if isEmptyString(contact.Id) {
		err += "Missing Id, "
	}
	if isEmptyString(contact.FirstName) {
		err += "Missing first name, "
	}
	if isEmptyString(contact.LastName) {
		err += "Missing last name, "
	}

	err = strings.TrimSuffix(err, ", ")

	return (len(err) == 0), fmt.Errorf(err)
}
