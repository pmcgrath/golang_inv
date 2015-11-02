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
	RabbitServers     xmlRabbitServers     `xml:"rabbitServers"`
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
	XMLName   xml.Name       `xml:"nlog"`
	Transform string         `xml:"Transform,attr"`
	Targets   xmlNLogTargets `xml:"targets"`
	Rules     xmlNLogRules   `xml:"rules"`
}

type xmlRabbitServers struct {
	XMLName   xml.Name           `xml:"rabbitServers"`
	Transform string             `xml:"Transform,attr"`
	Adds      []xmlAppSettingAdd `xml:"add"`
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
	Loggers   []xmlNLogRulesLogger `xml:"logger"`
}

type xmlNLogRulesLogger struct {
	XMLName   xml.Name `xml:"logger"`
	Name      string   `xml:"name,attr"`
	MinLevel  string   `xml:"minLevel,attr"`
	WriteTo   string   `xml:"writeTo,attr"`
	AppendTo  string   `xml:"appendTo,attr"`
	Transform string   `xml:"Transform,attr"`
}

type xmlRabbitServerAdd struct {
	XMLName   xml.Name `xml:"add"`
	Key       string   `xml:"key,attr"`
	Value     string   `xml:"value,attr"`
	Transform string   `xml:"Transform,attr"`
}

func parseConfigXmlContent(reader io.Reader) (xmlConfiguration, error) {
	var config xmlConfiguration
	if err := xml.NewDecoder(reader).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func mergeConfigXmlFileContents(base, update xmlConfiguration) xmlConfiguration {
	// This variable is just for clarity, as thei base value is passed by value - we could just mutate and return it
	merged := base

	// AppSettings
	if update.AppSettings.Transform == "Replace" {
		merged.AppSettings = update.AppSettings
	} else {
		// Deal with Inserts and SetAttributes
		for _, updateItem := range update.AppSettings.Adds {
			if updateItem.Transform == "Insert" {
				merged.AppSettings.Adds = append(merged.AppSettings.Adds, updateItem)
			} else if updateItem.Transform == "SetAttributes" {
				for index, mergedItem := range merged.AppSettings.Adds {
					if mergedItem.Key == updateItem.Key {
						merged.AppSettings.Adds[index].Value = updateItem.Value
					}
				}
			}
		}
	}

	// ConnectionStrings
	if update.ConnectionStrings.Transform == "Replace" {
		merged.ConnectionStrings = update.ConnectionStrings
	} else {
		// Deal with Inserts and SetAttributes
		for _, updateItem := range update.ConnectionStrings.Adds {
			if updateItem.Transform == "Insert" {
				merged.ConnectionStrings.Adds = append(merged.ConnectionStrings.Adds, updateItem)
			} else if updateItem.Transform == "SetAttributes" {
				for index, mergedItem := range merged.ConnectionStrings.Adds {
					if mergedItem.Name == updateItem.Name {
						// Only update if not empty
						if updateItem.ConnectionString != "" {
							merged.ConnectionStrings.Adds[index].ConnectionString = updateItem.ConnectionString
						}
						if updateItem.ProviderName != "" {
							merged.ConnectionStrings.Adds[index].ProviderName = updateItem.ProviderName
						}
					}
				}
			}
		}

	}

	// NLog
	if update.NLog.Transform == "Replace" {
		merged.NLog = update.NLog
	} else {
		// Targets
		if update.NLog.Targets.Transform == "Replace" {
			merged.NLog.Targets = update.NLog.Targets
		} else {
			// Deal with Inserts and SetAttributes
			for _, updateItem := range merged.NLog.Targets.Targets {
				if updateItem.Transform == "Insert" {
					merged.NLog.Targets.Targets = append(merged.NLog.Targets.Targets, updateItem)
				} else if updateItem.Transform == "SetAttributes" {
					for index, mergedItem := range merged.NLog.Targets.Targets {
						if mergedItem.Name == updateItem.Name {
							// Only update if not empty - revisit to use reflection
							if updateItem.Type != "" {
								merged.NLog.Targets.Targets[index].Type = updateItem.Type
							}
							if updateItem.Address != "" {
								merged.NLog.Targets.Targets[index].Address = updateItem.Address
							}
							if updateItem.Facility != "" {
								merged.NLog.Targets.Targets[index].Facility = updateItem.Facility
							}
							if updateItem.GelfServer != "" {
								merged.NLog.Targets.Targets[index].GelfServer = updateItem.GelfServer
							}
							if updateItem.Port != "" {
								merged.NLog.Targets.Targets[index].Port = updateItem.Port
							}
							if updateItem.MaxChunkSize != "" {
								merged.NLog.Targets.Targets[index].MaxChunkSize = updateItem.MaxChunkSize
							}
							if updateItem.GrayLogVersion != "" {
								merged.NLog.Targets.Targets[index].GrayLogVersion = updateItem.GrayLogVersion
							}
						}
					}
				}
			}
		}
		// Rules
		if update.NLog.Rules.Transform == "Replace" {
			merged.NLog.Rules = update.NLog.Rules
		} else {
			// Deal with Inserts and SetAttributes
			for _, updateItem := range merged.NLog.Rules.Loggers {
				if updateItem.Transform == "Insert" {
					merged.NLog.Rules.Loggers = append(merged.NLog.Rules.Loggers, updateItem)
				} else if updateItem.Transform == "SetAttributes" {
					for index, mergedItem := range merged.NLog.Rules.Loggers {
						if mergedItem.Name == updateItem.Name {
							// Only update if not empty - revisit to use reflection
							if updateItem.MinLevel != "" {
								merged.NLog.Rules.Loggers[index].MinLevel = updateItem.MinLevel
							}
							if updateItem.WriteTo != "" {
								merged.NLog.Rules.Loggers[index].WriteTo = updateItem.WriteTo
							}
							if updateItem.AppendTo != "" {
								merged.NLog.Rules.Loggers[index].AppendTo = updateItem.AppendTo
							}
						}
					}
				}
			}

		}
	}

	// RabbitServers
	if update.RabbitServers.Transform == "Replace" {
		merged.RabbitServers = update.RabbitServers
	} else {
		// Deal with Inserts and SetAttributes
		for _, updateItem := range update.RabbitServers.Adds {
			if updateItem.Transform == "Insert" {
				merged.RabbitServers.Adds = append(merged.RabbitServers.Adds, updateItem)
			} else if updateItem.Transform == "SetAttributes" {
				for index, mergedItem := range merged.RabbitServers.Adds {
					if mergedItem.Key == updateItem.Key {
						merged.RabbitServers.Adds[index].Value = updateItem.Value
					}
				}
			}
		}
	}

	return merged
}
