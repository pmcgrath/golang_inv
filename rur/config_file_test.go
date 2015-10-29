package main

import (
	"bytes"
	"testing"
)

func TestParseConfigXmlContent(t *testing.T) {
	content := bytes.NewBufferString(configContentFromService2)

	config, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse result, error: %#v", err)
	}
	t.Logf("-> %#v", config.AppSettings.Adds)
	t.Logf("-> %#v", config.ConnectionStrings.Adds)
	t.Logf("-> %#v", config.NLog.Targets)
}

func TestParseConfigXmlContentForTransformation(t *testing.T) {
	content := bytes.NewBufferString(transformationConfigContentFromService)

	config, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse result, error: %#v", err)
	}
	t.Logf("-> %#v", config.AppSettings.Adds)
	t.Logf("-> %#v", config.ConnectionStrings.Adds)
	t.Logf("-> %#v", config.NLog.Targets)
}
