package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type serviceConfiguration struct {
	Name         string
	Environments map[string]configuration
}

func (s serviceConfiguration) String() string {
	res := s.Name + "\n"
	for env, config := range s.Environments {
		res += fmt.Sprintf("**%s\n%s\n", env, config)
	}

	return res
}

type configuration struct {
	AppSettings map[string]string
	Databases   []database
	Loggers     []logger
}

func (c configuration) String() string {
	res := "AppSettings\n"
	for key, value := range c.AppSettings {
		res += fmt.Sprintf("\t%s = %s\n", key, value)
	}

	res += "Databases\n"
	for _, database := range c.Databases {
		res += fmt.Sprintf("\t%s\n", database)
	}

	res += "Loggers\n"
	for _, logger := range c.Loggers {
		res += fmt.Sprintf("\t%s\n", logger)
	}

	return res
}

type database struct {
	Type                   string
	Host                   string
	Port                   int
	Name                   string
	UsesIntegratedSecurity bool
	ConnectionString       string
}

func (d database) String() string {
	return fmt.Sprintf("Type = %s, Host = %s, Port = %d, Name = %s, Integrated Security = %t, ConnectionString = %s",
		d.Type,
		d.Host,
		d.Port,
		d.Name,
		d.UsesIntegratedSecurity,
		d.ConnectionString)
}

type logger struct {
	Name        string
	Level       string
	Target      string
	Facility    string
	Destination string
}

func (l logger) String() string {
	return fmt.Sprintf("Name = %s, Level = %s, Target = %s, Facility = %s, Destination = %s",
		l.Name,
		l.Level,
		l.Target,
		l.Facility,
		l.Destination)
}

func parseForService(directoryPath string) (serviceConfiguration, error) {
	config := serviceConfiguration{
		Name:         path.Base(directoryPath),
		Environments: make(map[string]configuration),
	}

	configFilePaths, err := getFSBasedRepoConfigFilePaths(directoryPath)
	if err != nil {
		return config, err
	}

	baseConfigFilePath := ""
	for _, configFilePath := range configFilePaths {
		_, configFileName := path.Split(configFilePath)
		configFileName = strings.ToLower(configFileName)
		if configFileName == "web.config" || configFileName == "app.config" {
			baseConfigFilePath = configFilePath
		}
	}
	if baseConfigFilePath == "" {
		return config, fmt.Errorf("No base config file found, must have a web.config or app.config which is the base config file")
	}

	baseXmlConfig, err := parseConfigXmlFile(baseConfigFilePath)
	if err != nil {
		return config, err
	}
	baseConfig := transformXmlConfig(baseXmlConfig)
	config.Environments["dev"] = baseConfig

	for _, configFilePath := range configFilePaths {
		if configFilePath == baseConfigFilePath {
			continue
		}

		_, configFileName := path.Split(configFilePath)

		xmlConfig, err := parseConfigXmlFile(configFilePath)
		if err != nil {
			return config, err
		}

		mergedXmlConfig := mergeConfigXmlFileContents(baseXmlConfig, xmlConfig)
		envConfig := transformXmlConfig(mergedXmlConfig)

		// Add entry for file - assuming no clashes for file name casing
		env := strings.ToLower(configFileName)
		env = env[0:strings.Index(env, ".")]
		config.Environments[env] = envConfig
	}

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

func transformXmlConfig(xmlConfig xmlConfiguration) configuration {
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

		// This is simlistic - have not had to probe\inspect to determine if a db and if so what type of db - seems to be working at this stage so I'm done
		switch strings.ToLower(connectionString.Name) {
		case "eventstore":
			// This only works if the name is consistent
			databases = append(databases, parseEventStoreConnectionString(connectionString.ConnectionString))
			continue
		case "metric", "metrics":
			// This only works if the name is consistent
			databases = append(databases, parseInfluxDBConnectionString(connectionString.ConnectionString))
			continue
		case "rabbitmq":
			// This isn't a database
			continue
		default:
			log.Printf("Connection string which we do not know how to process: name is [%s] and provider name is [%s]\n", connectionString.Name, connectionString.ConnectionString)
		}
	}

	loggers := transformNLogXml(xmlConfig.NLog)

	return configuration{
		AppSettings: appSettings,
		Databases:   databases,
		Loggers:     loggers,
	}
}

func parseMsSqlConnectionString(value string) database {
	db := database{Type: "MSSQL", ConnectionString: value}

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
	db := database{Type: "EventStore", ConnectionString: value}

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
			port, _ := strconv.Atoi(value[portSeperatorIndex+1:])

			db.Host = host
			db.Port = port
		}
	}

	return db
}

func parseInfluxDBConnectionString(value string) database {
	db := database{Type: "InfluxDB", ConnectionString: value}

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
		case "port":
			port, _ := strconv.Atoi(value)
			db.Port = port
		case "database":
			db.Name = value
		}
	}

	return db
}

func transformNLogXml(nlog xmlNLog) []logger {
	loggers := make([]logger, 0)
	for _, nlogLogger := range nlog.Rules.Loggers {
		logger := logger{
			Name:  nlogLogger.Name,
			Level: nlogLogger.MinLevel,
		}

		logger.Target = nlogLogger.WriteTo
		if nlogLogger.AppendTo != "" {
			logger.Target = nlogLogger.AppendTo
		}

		for _, nlogTarget := range nlog.Targets.Targets {
			if nlogTarget.Name == logger.Target {
				logger.Facility = nlogTarget.Facility
				logger.Destination = nlogTarget.GelfServer
				if nlogTarget.Address != "" {
					// Case of local udp target
					logger.Destination = nlogTarget.Address
				}
				break
			}
		}

		loggers = append(loggers, logger)
	}

	return loggers
}
