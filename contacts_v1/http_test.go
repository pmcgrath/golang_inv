package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestRequestContextIsLoggedIn(t *testing.T) {
	spec := &Spec{t}

	c := &RequestContext{}

	if c.IsLoggedIn() {
		t.Fatal("Should not be logged in")
	}

	c.Session = &Session{}
	c.Session.UserName = "Ted"

	spec.Assert(c.IsLoggedIn(), "Should be logged in")
}

func TestRequestContextiGetLogMessagePrefix(t *testing.T) {
	spec := &Spec{t}

	c := &RequestContext{
		Id: "TheId",
	}

	spec.Assert(c.GetLogMessagePrefix() == "TheId  []", "Unexpected message")

	c.Session = &Session{}
	c.Session.Id = "SID"
	c.Session.UserName = "Ted"

	spec.Assert(c.GetLogMessagePrefix() == "TheId SID [Ted]", "Unexpected message")
}
