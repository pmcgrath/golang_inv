package main

import (
	"errors"
	"fmt"
	"strings"
)

func getIdsFromUrlPath(path, pathPrefix string) ([]string, error) {
	if !strings.HasPrefix(path, pathPrefix) {
		return nil, fmt.Errorf("Prefix %s not found", pathPrefix)
	}

	pathWithoutPrefix := strings.TrimPrefix(path, pathPrefix)
	if pathWithoutPrefix == pathPrefix {
		return nil, errors.New("No ids exist")
	}

	ids := strings.Split(pathWithoutPrefix, "/")
	return ids, nil
}
