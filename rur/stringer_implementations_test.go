/*
	See 	http://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
		http://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
*/
package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

var configStringToAvoidCompilerOptimsations string // See http://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

func BenchmarkConfigurationStringerUsingConcat(b *testing.B) {
	runConfigurationStringerBenchmark(
		b,
		func(config *configuration) string {
			return config.StringConcat()
		})
}

func BenchmarkConfigurationStringerUsingByteBuffer(b *testing.B) {
	runConfigurationStringerBenchmark(
		b,
		func(config *configuration) string {
			return config.StringBuffer()
		})
}

func (c configuration) StringConcat() string {
	var buffer bytes.Buffer

	buffer.WriteString("AppSettings\n")
	for key, value := range c.AppSettings {
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(" = ")
		buffer.WriteString(value)
		buffer.WriteString("\n")
	}

	buffer.WriteString("Databases\n")
	for _, database := range c.Databases {
		buffer.WriteString("\tType = ")
		buffer.WriteString(database.Type)
		buffer.WriteString(", Host = ")
		buffer.WriteString(database.Host)
		buffer.WriteString(", Name = ")
		buffer.WriteString(database.Name)
		buffer.WriteString(", Integrated Security = ")
		buffer.WriteString(strconv.FormatBool(database.UsesIntegratedSecurity))
		buffer.WriteString("\n")
	}

	buffer.WriteString("Loggers\n")
	for _, logger := range c.Loggers {
		buffer.WriteString("\tName = ")
		buffer.WriteString(logger.Name)
		buffer.WriteString(", Level = ")
		buffer.WriteString(logger.Level)
		buffer.WriteString(", Facility = ")
		buffer.WriteString(logger.Facility)
		buffer.WriteString(", Destination = ")
		buffer.WriteString(logger.Destination)
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (c configuration) StringBuffer() string {
	res := "AppSettings\n"
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

	res += "Loggers\n"
	for _, logger := range c.Loggers {
		res += fmt.Sprintf("\tName = %s, Level = %s, Facility = %s, Destination = %s\n",
			logger.Name,
			logger.Level,
			logger.Facility,
			logger.Destination)
	}

	return res
}

func runConfigurationStringerBenchmark(b *testing.B, stringerFuncWrapper func(*configuration) string) {
	config := &configuration{
		AppSettings: map[string]string{
			"Key1": "Value1",
			"Key2": "Value2",
		},
		Databases: []database{
			{Type: "MSSQL", Host: "Source1", Name: "Database1", UsesIntegratedSecurity: false},
			{Type: "EventStore", Host: "Source2", Name: "", UsesIntegratedSecurity: true},
		},
		Loggers: []logger{
			{Name: "Name1", Level: "INFO", Facility: "Facility1", Destination: "Destination1"},
			{Name: "Name1", Level: "DEBUG", Facility: "Facility1", Destination: "Destination1"},
		}}

	// Capture result to avoid compiler optimisation kicking in
	// http://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
	var configString string
	for n := 0; n < b.N; n++ {
		configString = stringerFuncWrapper(config)
	}
	configStringToAvoidCompilerOptimsations = configString
}
