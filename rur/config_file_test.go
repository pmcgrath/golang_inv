package main

import (
	"bytes"
	"testing"
)

func TestParseConfigXmlContent(t *testing.T) {
	content := bytes.NewBufferString(configContentSample)

	config, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse result, error: %#v", err)
	}
	// AppSettings
	appSettings := config.AppSettings.Adds
	if len(appSettings) != 3 {
		t.Errorf("Unexpected appSettings count, expected %d but got %d", 3, len(appSettings))
	}
	if !assertAppSettingMatch(appSettings[0], "IncludeJsonTypeInfo", "true", "") {
		t.Errorf("Unexpected appSettings[0] %#v", appSettings[0])
	}
	if !assertAppSettingMatch(appSettings[1], "Metrics.GlobalContextName", "ServiceABC", "") {
		t.Errorf("Unexpected appSettings[1] %#v", appSettings[1])
	}
	if !assertAppSettingMatch(appSettings[2], "SecurityServiceUrl", "http://dev.security.com/", "") {
		t.Errorf("Unexpected appSettings[2] %#v", appSettings[2])
	}
	// ConnectionStrings
	connectionStrings := config.ConnectionStrings.Adds
	if len(connectionStrings) != 4 {
		t.Errorf("Unexpected connectionStrings count, expected %d but got %d", 3, len(connectionStrings))
	}
	if !assertConnectionStringMatch(connectionStrings[0], "AAADatabase", "Data Source=db1; Initial Catalog=AAA; MultiSubnetFailover=True; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[0] %#v", connectionStrings[0])
	}
	if !assertConnectionStringMatch(connectionStrings[1], "BBBDatabase", "Data Source=db1; Initial Catalog=BBB; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[1] %#v", connectionStrings[1])
	}
	if !assertConnectionStringMatch(connectionStrings[2], "CCDatabase", "Data Source=db2; Initial Catalog=CCC; MultiSubnetFailover=True; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[2] %#v", connectionStrings[2])
	}
	if !assertConnectionStringMatch(connectionStrings[3], "EventStore", "CONNECTTo=tcp://admin:ASuperDupperStrongPassword@127.0.0.1:1113", "", "") {
		t.Errorf("Unexpected connectionStrings[3] %#v", connectionStrings[3])
	}
	// Nlog
	logTargets := config.NLog.Targets.Targets
	if len(logTargets) != 2 {
		t.Errorf("Unexpected log targets count, expected %d but got %d", 2, len(logTargets))
	}
	if !assertLogTargetMatch(logTargets[0], "UdpOutlet", "NLogViewer", "udp://localhost:7071", "", "", "", "", "", "") {
		t.Errorf("Unexpected logTargets[0] %#v", logTargets[0])
	}
	if !assertLogTargetMatch(logTargets[1], "Gelf", "Gelf", "", "log.company.com", "ServiceABC", "12201", "8154", "0.9.6", "") {
		t.Errorf("Unexpected logTargets[1] %#v", logTargets[1])
	}
	logLoggers := config.NLog.Rules.Loggers
	if len(logLoggers) != 2 {
		t.Errorf("Unexpected logLoggers count, expected %d but got %d", 2, len(logLoggers))
	}
	if !assertLogLoggerMatch(logLoggers[0], "*", "Trace", "UdpOutlet", "", "") {
		t.Errorf("Unexpected logLoggers[0] %#v", logLoggers[0])
	}
	if !assertLogLoggerMatch(logLoggers[1], "*", "Trace", "", "Gelf", "") {
		t.Errorf("Unexpected logLoggers1] %#v", logLoggers[1])
	}
}

func TestParseConfigXmlContentForTransformationFile(t *testing.T) {
	content := bytes.NewBufferString(configContentTransformation)

	config, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse result, error: %#v", err)
	}
	// AppSettings
	appSettings := config.AppSettings.Adds
	if len(appSettings) != 3 {
		t.Errorf("Unexpected appSettings count, expected %d but got %d", 3, len(appSettings))
	}
	if !assertAppSettingMatch(appSettings[0], "IsLive", "True", "Insert") {
		t.Errorf("Unexpected appSettings[0] %#v", appSettings[0])
	}
	if !assertAppSettingMatch(appSettings[1], "IncludeJsonTypeInfo", "false", "SetAttributes") {
		t.Errorf("Unexpected appSettings[1] %#v", appSettings[1])
	}
	if !assertAppSettingMatch(appSettings[2], "SecurityServiceUrl", "https://security.services.local/api2/security/", "SetAttributes") {
		t.Errorf("Unexpected appSettings[2] %#v", appSettings[2])
	}
	// ConnectionStrings
	connectionStrings := config.ConnectionStrings.Adds
	if len(connectionStrings) != 4 {
		t.Errorf("Unexpected connectionStrings count, expected %d but got %d", 4, len(connectionStrings))
	}
	if !assertConnectionStringMatch(connectionStrings[0], "AAADatabase", "Server=tcp:db2; Database=AAA; MultiSubnetFailover=True; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[0] %#v", connectionStrings[0])
	}
	if !assertConnectionStringMatch(connectionStrings[1], "BBBDatabase", "Data Source=db2;Initial Catalog=BBB; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[1] %#v", connectionStrings[1])
	}
	if !assertConnectionStringMatch(connectionStrings[2], "CCCDatabase", "Data Source=cdb2; Initial Catalog=CCC; Integrated Security=SSPI;", "System.Data.SqlClient", "") {
		t.Errorf("Unexpected connectionStrings[2] %#v", connectionStrings[2])
	}
	if !assertConnectionStringMatch(connectionStrings[3], "metrics", "host=metrics01;port=8086;user=root;password=PASSWORD;database=metrics", "", "") {
		t.Errorf("Unexpected connectionStrings[3] %#v", connectionStrings[3])
	}
	// Nlog
	logTargets := config.NLog.Targets.Targets
	if len(logTargets) != 1 {
		t.Errorf("Unexpected log targets count, expected %d but got %d", 2, len(logTargets))
	}
	if !assertLogTargetMatch(logTargets[0], "Gelf", "", "", "log.company.com", "", "", "", "", "SetAttributes") {
		t.Errorf("Unexpected logTargets[0] %#v", logTargets[0])
	}
	logLoggers := config.NLog.Rules.Loggers
	if len(logLoggers) != 2 {
		t.Errorf("Unexpected logLoggers count, expected %d but got %d", 2, len(logLoggers))
	}
	if !assertLogLoggerMatch(logLoggers[0], "*", "Warn", "", "GelfOld", "") {
		t.Errorf("Unexpected logLoggers[0] %#v", logLoggers[0])
	}
	if !assertLogLoggerMatch(logLoggers[1], "*", "Debug", "", "Gelf", "") {
		t.Errorf("Unexpected logLoggers[1] %#v", logLoggers[1])
	}
}

func TestMergeConfigXmlFileContents(t *testing.T) {
	content := bytes.NewBufferString(configContentSample)
	base, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse error for base : %#v", err)
	}

	content = bytes.NewBufferString(configContentTransformation)
	update, err := parseConfigXmlContent(content)
	if err != nil {
		t.Errorf("Unexpected parse error for update : %#v", err)
	}

	t.Logf("\n\n**** base before\n%s", transformXmlConfig(base))
	//	update.AppSettings.Transform = "Replace"
	update.ConnectionStrings.Transform = "Replace"
	merged := mergeConfigXmlFileContents(base, update)

	t.Logf("\n\nbase after\n%s", transformXmlConfig(base))
	t.Logf("\n\nupdate\n%s", transformXmlConfig(update))
	t.Logf("\n\nmerged\n%s", transformXmlConfig(merged))
}

func assertAppSettingMatch(actual xmlAppSettingAdd, key, value, transform string) bool {
	return actual.Key == key &&
		actual.Value == value &&
		actual.Transform == transform
}

func assertConnectionStringMatch(actual xmlConnectionStringAdd, name, connectionString, providerName, transform string) bool {
	return actual.Name == name &&
		actual.ConnectionString == connectionString &&
		actual.ProviderName == providerName &&
		actual.Transform == transform
}

func assertLogTargetMatch(actual xmlNLogTargetsTarget, name, theType, address, gelfServer, facility, port, maxChunkSize, grayLogVersion, transform string) bool {
	return actual.Name == name &&
		actual.Type == theType &&
		actual.Address == address &&
		actual.GelfServer == gelfServer &&
		actual.Facility == facility &&
		actual.Port == port &&
		actual.MaxChunkSize == maxChunkSize &&
		actual.GrayLogVersion == grayLogVersion &&
		actual.Transform == transform
}

func assertLogLoggerMatch(actual xmlNLogRulesLogger, name, minLevel, writeTo, appendTo, transform string) bool {
	return actual.Name == name &&
		actual.MinLevel == minLevel &&
		actual.WriteTo == writeTo &&
		actual.AppendTo == appendTo &&
		actual.Transform == transform
}
