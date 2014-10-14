package main

import (
	"testing"
)

func TestIsEmptyString(t *testing.T) {
	spec := &Spec{t}

	testCases := []struct {
		s        string // Input
		expected bool   // Expected result
	}{
		{"", true},
		{" ", true},
		{"     ", true},
		{"1", false},
		{"aaaa", false},
	}

	for _, testCase := range testCases {
		actual := isEmptyString(testCase.s)
		spec.Assert(actual == testCase.expected, "Unexpected result %t for input [%s]", actual, testCase.s)
	}
}
