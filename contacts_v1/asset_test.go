package main

import (
	"testing"
)

func TestGetAssetContentType(t *testing.T) {
	spec := &Spec{t}

	testCases := []struct {
		path     string
		expected string
	}{
		{"/assets/js/ted.js", "application/javascript"},
		{"/assets/css/a.css", "text/css; charset=utf-8"},
	}

	for _, testCase := range testCases {
		actual := getAssetContentType(testCase.path)
		spec.Assert(actual == testCase.expected, "Unexpected result %s for input [%s]", actual, testCase.path)
	}
}
