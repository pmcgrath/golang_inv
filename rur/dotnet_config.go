package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type configuration struct {
	ServiceName string
	AppSettings map[string]string
	Databases   []database
	LogTargets  []logTarget
}

func (c configuration) String() string {
	res := c.ServiceName + "\n"

	res += "AppSettings\n"
	for key, value := range c.AppSettings {
		res += fmt.Sprintf("\t%s = %s\n", key, value)
	}

	res += "Databases\n"
	for _, database := range c.Databases {
		res += fmt.Sprintf("\tType = %s, Host = %s, Name = %s, Integrated Security = %t\n",
			database.Type,
			database.Host,
			database.Name,
			database.UsesIntegratedSecurity)
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

type database struct {
	Type                   string
	Host                   string
	Name                   string
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

	databases := make([]database, 0)
	for _, connectionString := range xmlConfig.ConnectionStrings.Adds {
		if strings.ToLower(connectionString.ProviderName) == "system.data.sqlclient" {
			databases = append(databases, parseMsSqlConnectionString(connectionString.ConnectionString))
			continue
		}

		// This is simlistic - have not had to probe\inspacet to determine if a db and if so what type of db - seems to be working at this stage so I'm done
		switch strings.ToLower(connectionString.Name) {
		case "eventstore":
			// This only works if the name is consistent
			databases = append(databases, parseEventStoreConnectionString(connectionString.ConnectionString))
			continue
		case "metrics":
			// This only works if the name is consistent
			databases = append(databases, parseInfluxDBConnectionString(connectionString.ConnectionString))
			continue
		case "rabbitmq":
			// This isn't a database
			continue
		default:
			log.Printf("%s : Connection string which we do not know how to process: name is [%s] and provider name is [%s]\n", serviceName, connectionString.Name, connectionString.ConnectionString)
		}
	}

	logTargets := transformNLogXml(xmlConfig.NLog)

	return configuration{
		ServiceName: serviceName,
		AppSettings: appSettings,
		Databases:   databases,
		LogTargets:  logTargets,
	}
}

func parseMsSqlConnectionString(value string) database {
	db := database{Type: "MSSQL"}

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
			db.Name = value
		case "integrated security":
			if strings.ToLower(value) == "sspi" {
				db.UsesIntegratedSecurity = true
			}
		}
	}

	return db
}

func parseEventStoreConnectionString(value string) database {
	db := database{Type: "EventStore"}

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
		case "connectto":
			// Format is protocol://username:password@host:port i.e. tcp://admin:ASuperDupperStrongPassword@127.0.0.1:1113 cluster://admin:ASuperDupperStrongPassword@eventstore.dev.local:1113
			hostSeperatorIndex := strings.Index(value, "@")
			portSeperatorIndex := strings.LastIndex(value, ":") // Is port optional ? Use a default if not supplied
			host := value[hostSeperatorIndex+1 : portSeperatorIndex]

			db.Host = host
		}
	}

	return db
}

func parseInfluxDBConnectionString(value string) database {
	db := database{Type: "InfluxDB"}

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
		case "host":
			db.Host = value
		case "database":
			db.Name = value
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
