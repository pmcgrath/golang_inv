package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type configuration struct {
	ServiceName    string
	AppSettings    map[string]string
	MsSqlDatabases []msSqlDatabase
	LogTargets     []logTarget
}

func (c configuration) String() string {
	res := c.ServiceName + "\n"

	res += "AppSettings\n"
	for key, value := range c.AppSettings {
		res += fmt.Sprintf("\t%s = %s\n", key, value)
	}

	res += "MsSqlDatabases\n"
	for _, msSqlDatabase := range c.MsSqlDatabases {
		res += fmt.Sprintf("\tSource = %s, Database = %s, Integrated Security = %t\n",
			msSqlDatabase.Source,
			msSqlDatabase.Database,
			msSqlDatabase.UsesIntegratedSecurity)
	}

	res += "LogTargets\n"
	for _, logTarget := range c.LogTargets {
		res += fmt.Sprintf("\tName = %s, Facility = %s, Destination = %s\n",
			logTarget.Name,
			logTarget.Facility,
			logTarget.Destination)
	}

	return res
}

type msSqlDatabase struct {
	Source                 string
	Database               string
	UsesIntegratedSecurity bool
}

type logTarget struct {
	Name        string
	Facility    string
	Destination string
}

func parseForService(directoryPath string) (configuration, error) {
	serviceName := path.Base(directoryPath)

	mainProjectDirectoryPath := path.Join(directoryPath, serviceName)
	webConfigFilePath := path.Join(mainProjectDirectoryPath, "web.config")
	webConfigFileExists := testIfFileExists(webConfigFilePath)
	if !webConfigFileExists {
		return configuration{}, fmt.Errorf("No web.config file exists for service %s", serviceName)
	}

	webConfigFile, err := os.Open(webConfigFilePath)
	if err != nil {
		return configuration{}, err
	}
	defer webConfigFile.Close()

	xmlConfig, err := parseConfigXmlContent(webConfigFile)
	if err != nil {
		return configuration{}, err
	}

	config := transformXmlConfig(serviceName, xmlConfig)

	return config, nil
}

func transformXmlConfig(serviceName string, xmlConfig xmlConfiguration) configuration {
	appSettings := make(map[string]string, len(xmlConfig.AppSettings.Adds))
	for _, appSetting := range xmlConfig.AppSettings.Adds {
		appSettings[appSetting.Key] = appSetting.Value
	}

	msSqlDatabases := make([]msSqlDatabase, 0)
	for _, connectionString := range xmlConfig.ConnectionStrings.Adds {
		if connectionString.ProviderName == "" {
			log.Printf("%s : Connection string with no provider name : %s\n", serviceName, connectionString.ConnectionString)
		}
		if connectionString.ProviderName == "System.Data.SqlClient" {
			msSqlDatabases = append(msSqlDatabases, parseMsSqlConnectionString(connectionString.ConnectionString))
		}
	}

	logTargets := transformNLogXml(xmlConfig.NLog)

	return configuration{
		ServiceName:    serviceName,
		AppSettings:    appSettings,
		MsSqlDatabases: msSqlDatabases,
		LogTargets:     logTargets,
	}
}

func parseMsSqlConnectionString(value string) msSqlDatabase {
	db := msSqlDatabase{}

	attributes := strings.Split(value, ";")
	for _, attribute := range attributes {
		// Cater for trailing ; in which case we will have an empty attribute
		if strings.TrimSpace(attribute) == "" {
			continue
		}

		attributeParts := strings.Split(attribute, "=")
		key := strings.TrimSpace(strings.ToLower(attributeParts[0]))
		value := strings.TrimSpace(attributeParts[1])

		switch key {
		case "data source", "server":
			db.Source = value
		case "database", "initial catalog":
			db.Database = value
		case "integrated security":
			if strings.ToLower(value) == "sspi" {
				db.UsesIntegratedSecurity = true
			}
		}
	}

	return db
}

func transformNLogXml(nlog xmlNLog) []logTarget {
	logTargets := make([]logTarget, 0)
	for _, nlogRule := range nlog.Rules.Rules {
		logTarget := logTarget{}

		logTarget.Name = nlogRule.WriteTo
		if nlogRule.AppendTo != "" {
			logTarget.Name = nlogRule.AppendTo
		}

		for _, nlogTarget := range nlog.Targets.Targets {
			if nlogTarget.Name == logTarget.Name {
				logTarget.Facility = nlogTarget.Facility
				logTarget.Destination = nlogTarget.GelfServer
				if nlogTarget.Address != "" {
					// Case of local udp target
					logTarget.Destination = nlogTarget.Address
				}
				break
			}
		}

		logTargets = append(logTargets, logTarget)
	}

	return logTargets
}
