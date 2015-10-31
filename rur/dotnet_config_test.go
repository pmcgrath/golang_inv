package main

import "testing"

func TestConfigurationString(t *testing.T) {
	config := configuration{
		ServiceName: "Ted",
		AppSettings: map[string]string{
			"Key1": "Value1",
			"Key2": "value2"},
		Databases: []database{
			database{Type: "SuperDB1", Host: "MyDBServer", Port: 1433, Name: "MyDB1", UsesIntegratedSecurity: false, ConnectionString: "connect to MyDB1"},
			database{Type: "SuperDB2", Host: "MyDBServer", Port: 0, Name: "MyDB2", UsesIntegratedSecurity: false, ConnectionString: "connect to MyDB2"},
		},
		LogTargets: []logTarget{
			logTarget{Name: "log", Facility: "myservice", Destination: "logs.service.com"},
		},
	}

	value := config.String()

	t.Logf(value)
}

func TestParseMsSqlConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString               string
		ExpectedHost                   string
		ExpectedName                   string
		ExpectedUsesIntegratedSecurity bool
	}{
		{"Server=tcp:dbserver1; Database=myDB; MultiSubnetFailover=True; Integrated Security=SSPI;", "tcp:dbserver1", "myDB", true},
		{"Data Source=AAA; Initial Catalog=BBB; uid=auser; password=itsasecret", "AAA", "BBB", false},
		{"Server=DBServer", "DBServer", "", false},
		{"Server=DBServer;           Initial catalog    =    DB1", "DBServer", "DB1", false},
		{";", "", "", false},
	} {
		db := parseMsSqlConnectionString(testCase.ConnectionString)

		if db.Host != testCase.ExpectedHost {
			t.Errorf("Unexpected host, expected : %s actual : %s", testCase.ExpectedHost, db.Host)
		}
		if db.Name != testCase.ExpectedName {
			t.Errorf("Unexpected name, expected : %s actual : %s", testCase.ExpectedName, db.Name)
		}
		if db.UsesIntegratedSecurity != testCase.ExpectedUsesIntegratedSecurity {
			t.Errorf("Unexpected uses integrated security: %t", testCase.ExpectedUsesIntegratedSecurity, db.UsesIntegratedSecurity)
		}
		if db.ConnectionString != testCase.ConnectionString {
			t.Errorf("Unexpected connection string, expected : %s actual : %s", testCase.ConnectionString, db.ConnectionString)
		}
	}
}

func TestParseEventStoreConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString string
		ExpectedHost     string
		ExpectedPort     int
	}{
		{"CONNECTTo=tcp://admin:ASuperDupperStrongPassword@127.0.0.1:1113", "127.0.0.1", 1113},
		{"ConnectTo = cluster://admin:ASuperDupperStrongPassword@eventstore.dev.local:1113", "eventstore.dev.local", 1113},
	} {
		db := parseEventStoreConnectionString(testCase.ConnectionString)

		if db.Host != testCase.ExpectedHost {
			t.Errorf("Unexpected host, expected : %s actual : %s", testCase.ExpectedHost, db.Host)
		}
		if db.Port != testCase.ExpectedPort {
			t.Errorf("Unexpected port, expected : %d actual : %d", testCase.ExpectedPort, db.Port)
		}
		if db.Name != "" {
			t.Errorf("Unexpected name, expected empty, actual : %s", db.Name)
		}
		if db.UsesIntegratedSecurity {
			t.Errorf("Unexpected uses integrated security, expected false but got true")
		}
		if db.ConnectionString != testCase.ConnectionString {
			t.Errorf("Unexpected connection string, expected : %s actual : %s", testCase.ConnectionString, db.ConnectionString)
		}
	}
}

func TestParseInfluxDBConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString string
		ExpectedHost     string
		ExpectedPort     int
		ExpectedName     string
	}{
		{"host=metrics01;port=8086;user=root;password=PASSWORD;database=metrics", "metrics01", 8086, "metrics"},
	} {
		db := parseInfluxDBConnectionString(testCase.ConnectionString)

		if db.Host != testCase.ExpectedHost {
			t.Errorf("Unexpected host, expected : %s actual : %s", testCase.ExpectedHost, db.Host)
		}
		if db.Port != testCase.ExpectedPort {
			t.Errorf("Unexpected port, expected : %d actual : %d", testCase.ExpectedPort, db.Port)
		}
		if db.Name != testCase.ExpectedName {
			t.Errorf("Unexpected name, expected : %s actual : %s", testCase.ExpectedName, db.Name)
		}
		if db.UsesIntegratedSecurity {
			t.Errorf("Unexpected uses integrated security, expected false but got true")
		}
		if db.ConnectionString != testCase.ConnectionString {
			t.Errorf("Unexpected connection string, expected : %s actual : %s", testCase.ConnectionString, db.ConnectionString)
		}
	}
}

func TestTransformNLogXml(t *testing.T) {
	source := xmlNLog{
		Targets: xmlNLogTargets{
			Targets: []xmlNLogTargetsTarget{
				{
					Name:    "udpLogger",
					Type:    "NlogViewer",
					Address: "udp://localhost:7071",
				},
				{
					Name:       "gelfLogger",
					Type:       "Gelf",
					Facility:   "super_service",
					GelfServer: "log.company.com",
					Port:       "12201",
				},
				{
					Name:       "orphanedGelfLogger",
					Type:       "Gelf",
					Facility:   "super_service",
					GelfServer: "old.log.company.com",
					Port:       "12201",
				},
			},
		},
		Rules: xmlNLogRules{
			Rules: []xmlNLogRulesLogger{
				{
					Name:     "*",
					MinLevel: "Trace",
					WriteTo:  "udpLogger",
				},
				{
					Name:     "*",
					MinLevel: "Trace",
					AppendTo: "gelfLogger",
				},
			},
		},
	}

	logTargets := transformNLogXml(source)

	if len(logTargets) != 2 {
		t.Errorf("Unexpected log targets count, expected : 2 actual : %d", len(logTargets))
	}

	t.Logf("---> %#v\n\n", logTargets)
}
