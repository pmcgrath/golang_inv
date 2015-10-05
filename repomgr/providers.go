package main

import (
	"fmt"
	"strings"
)

type newProviderFn func(providerConnectionAttributes) (provider, error)

func getNewProviderFns() map[string]newProviderFn {
	return map[string]newProviderFn{
		"github": func(c providerConnectionAttributes) (provider, error) { return newGitHubProvider(c), nil },
		"stash":  func(c providerConnectionAttributes) (provider, error) { return newStashProvider(c), nil },
	}
}

func newProvider(providerName string, connAttrs providerConnectionAttributes) (provider, error) {
	newProviderFns := getNewProviderFns()
	newProviderFn, ok := newProviderFns[strings.ToLower(providerName)]
	if !ok {
		return nil, fmt.Errorf("Unknown SCM provider [%s]\n", providerName)
	}

	logDebugf("About to construct provider [%s]\n", providerName)
	provider, err := newProviderFn(connAttrs)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
