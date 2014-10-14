package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestGetOrDefaultEnvWhereFound(t *testing.T) {
	spec := &Spec{t}

	spec.Assert(GetOrDefaultEnv("PATH", ".") != ".", "PATH value not found")
}

func TestGetOrDefaultEnvWhereNotFound(t *testing.T) {
	spec := &Spec{t}

	spec.Assert(GetOrDefaultEnv("DOESNOTEXIST", ".") == ".", "DOESNOTEXIST value found")
}
