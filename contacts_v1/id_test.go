package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestUuid(t *testing.T) {
	spec := &Spec{t}

	id := Uuid()

	spec.Assert(len(id) == 36, "Unexpected id %s", id)
}

func TestUuidMultiples(t *testing.T) {
	spec := &Spec{t}

	idsChan := make(chan string, 50)
	genIdFunc := func(result chan<- string) { result <- Uuid() }

	for i := 1; i <= 50; i++ {
		genIdFunc(idsChan)
	}

	generatedIds := make(map[string]interface{})
	for i := 1; i <= 50; i++ {
		id := <-idsChan
		if _, ok := generatedIds[id]; ok {
			spec.Assert(false, "Duplicate id detected")
		}

		generatedIds[id] = struct{}{}
	}
}
