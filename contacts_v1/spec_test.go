package main

import (
	"testing"
)

type Spec struct {
	*testing.T
}

func (s *Spec) Assert(assertionResult bool, message string, messageArguments ...interface{}) {
	if !assertionResult {
		if messageArguments == nil {
			s.Fatal(message)
		} else {
			s.Fatalf(message, messageArguments...)
		}
	}
}
