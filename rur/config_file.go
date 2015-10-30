/*
	See 	http://www.goinggo.net/2013/06/reading-xml-documents-in-go.html
		https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/07.1.html
*/
package main

import (
	"encoding/xml"
	"io"
)

// This is a union of every setting that we are interested in, if we are not interested I do not include
type xmlConfiguration struct {
	XMLName           xml.Name             `xml:"configuration"`
	AppSettings       xmlAppSettings       `xml:"appSettings"`
	ConnectionStrings xmlConnectionStrings `xml:"connectionStrings"`
	NLog              xmlNLog              `xml:"nlog"`
}

type xmlAppSettings struct {
	XMLName   xml.Name           `xml:"appSettings"`
	Transform string             `xml:"Transform,attr"`
	Adds      []xmlAppSettingAdd `xml:"add"`
}

type xmlConnectionStrings struct {
	XMLName   xml.Name                 `xml:"connectionStrings"`
	Transform string                   `xml:"Transform,attr"`
	Adds      []xmlConnectionStringAdd `xml:"add"`
}

type xmlNLog struct {
	XMLName xml.Name       `xml:"nlog"`
	Targets xmlNLogTargets `xml:"targets"`
	Rules   xmlNLogRules   `xml:"rules"`
}

type xmlAppSettingAdd struct {
	XMLName   xml.Name `xml:"add"`
	Key       string   `xml:"key,attr"`
	Value     string   `xml:"value,attr"`
	Transform string   `xml:"Transform,attr"`
}

type xmlConnectionStringAdd struct {
	XMLName          xml.Name `xml:"add"`
	Name             string   `xml:"name,attr"`
	ConnectionString string   `xml:"connectionString,attr"`
	ProviderName     string   `xml:"providerName,attr"`
	Transform        string   `xml:"Transform,attr"`
}

type xmlNLogTargets struct {
	XMLName   xml.Name               `xml:"targets"`
	Transform string                 `xml:"Transform,attr"`
	Targets   []xmlNLogTargetsTarget `xml:"target"`
}

type xmlNLogTargetsTarget struct {
	XMLName        xml.Name `xml:"target"`
	Name           string   `xml:"name,attr"`
	Type           string   `xml:"type,attr"`
	Address        string   `xml:"address,attr"`
	Facility       string   `xml:"facility,attr"`
	GelfServer     string   `xml:"gelfserver,attr"`
	Port           string   `xml:"port,attr"`
	MaxChunkSize   string   `xml:"maxchunksize,attr"`
	GrayLogVersion string   `xml:"graylogversion,attr"`
	Transform      string   `xml:"Transform,attr"`
}

type xmlNLogRules struct {
	XMLName   xml.Name             `xml:"rules"`
	Transform string               `xml:"Transform,attr"`
	Rules     []xmlNLogRulesLogger `xml:"logger"`
}

type xmlNLogRulesLogger struct {
	XMLName   xml.Name `xml:"logger"`
	Name      string   `xml:"name,attr"`
	MinLevel  string   `xml:"minLevel,attr"`
	WriteTo   string   `xml:"writeTo,attr"`
	AppendTo  string   `xml:"appendTo,attr"`
	Transform string   `xml:"Transform,attr"`
}

func parseConfigXmlContent(reader io.Reader) (xmlConfiguration, error) {
	var config xmlConfiguration
	if err := xml.NewDecoder(reader).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
