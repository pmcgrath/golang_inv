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
		res += fmt.Sprintf("\tHost = %s, Database = %s, Integrated Security = %t\n",
			msSqlDatabase.Host,
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
	Host                   string
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
	configFilePaths, err := getFSBasedRepoConfigFilePaths(directoryPath)
	if err != nil {
		return configuration{}, err
	}

	xmlConfigs := make(map[string]xmlConfiguration)
	for _, configFilePath := range configFilePaths {
		_, configFileName := path.Split(configFilePath)

		xmlConfig, err := parseConfigXmlFile(configFilePath)
		if err != nil {
			return configuration{}, err
		}

		// Add entry for file - assuming no clashes for file name casing
		key := strings.ToLower(configFileName)
		xmlConfigs[key] = xmlConfig

		transformXmlConfig(serviceName, xmlConfig)
	}

	// pmcg HACK for now just deal with web.config
	config := transformXmlConfig(serviceName, xmlConfigs["web.config"])

	return config, nil
}

func parseConfigXmlFile(filePath string) (xmlConfiguration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return xmlConfiguration{}, err
	}
	defer file.Close()

	xmlConfig, err := parseConfigXmlContent(file)
	if err != nil {
		return xmlConfig, err
	}

	return xmlConfig, nil
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
			db.Host = value
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
