package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestGetIdsFromUrlPathSuccess(t *testing.T) {
	spec := &Spec{t}

	path := "/contacts/ted/toe"
	pathPrefix := "/contacts/"

	ids, err := getIdsFromUrlPath(path, pathPrefix)

	spec.Assert(err == nil, "Unexpected error %s", err)
	spec.Assert(len(ids) == 2, "Unexpected ids %s", ids)
	spec.Assert(ids[0] == "ted", "Unexpected id[0] of %s", ids[0])
	spec.Assert(ids[1] == "toe", "Unexpected id[1] of %s", ids[1])
}

func TestGetIdsFromUrlPathPrefixNotPresent(t *testing.T) {
	spec := &Spec{t}

	path := "/contacts/ted/toe"
	pathPrefix := "/notthere/"

	_, err := getIdsFromUrlPath(path, pathPrefix)

	spec.Assert(err != nil, "Expected error but none found")
	spec.Assert(err.Error() == "Prefix /notthere/ not found", "Unexpected error", err)
}
