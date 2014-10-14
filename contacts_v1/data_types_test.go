package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestUserAuthentication(t *testing.T) {
	spec := &Spec{t}

	user := &User{
		Password: "Tim",
	}

	spec.Assert(user.Authenticate("Bad") == false, "Should have failed")
	spec.Assert(user.Authenticate("Tim"), "Should have passed")
}

func TestUserGetContactIndexWhereContactExists(t *testing.T) {
	spec := &Spec{t}

	user := &User{
		Id: "pmcgrath",
		Contacts: []Contact{
			Contact{
				Id: "c1",
			},
			Contact{
				Id: "c2",
			},
			Contact{
				Id: "c3",
			},
		},
	}

	index, ok := user.GetContactIndex("c3")

	spec.Assert(ok, "Contact not found")
	spec.Assert(index == 2, "Unexpected index %d", index)
}

func TestUserGetContactIndexWhereContactiDoesNotExist(t *testing.T) {
	spec := &Spec{t}

	user := &User{}

	index, ok := user.GetContactIndex("c3")

	spec.Assert(ok == false, "Contact found")
	spec.Assert(index == -1, "Unexpected index %d", index)
}

func TestContactIsValidForSavingForValidCase(t *testing.T) {
	spec := &Spec{t}

	testCases := []struct {
		c             *Contact // Input
		expected      bool     // Expected result
		expectedError string   // Expected error string
	}{
		{c: &Contact{Id: "Id1", FirstName: "Ted", LastName: "Toe"}, expected: true, expectedError: ""},
		{c: &Contact{Id: "Id1", LastName: "Toe"}, expected: false, expectedError: "Missing first name"},
		{c: &Contact{Id: "Id1"}, expected: false, expectedError: "Missing first name, Missing last name"},
	}

	for _, testCase := range testCases {
		actual, actualError := testCase.c.IsValidForSaving()
		spec.Assert(actual == testCase.expected, "Unexpected result %t for input [%v]", actual, testCase.c)
		spec.Assert(actualError.Error() == testCase.expectedError, "Unexpected error %s for input [%v]", actualError.Error(), testCase.c)
	}
}
