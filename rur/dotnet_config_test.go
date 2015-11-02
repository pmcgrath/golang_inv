package main

import "testing"

func TestConfigurationString(t *testing.T) {
	config := configuration{
		AppSettings: map[string]string{
			"Key1": "Value1",
			"Key2": "value2"},
		Databases: []database{
			database{Type: "SuperDB1", Host: "MyDBServer", Port: 1433, Name: "MyDB1", UsesIntegratedSecurity: false, ConnectionString: "connect to MyDB1"},
			database{Type: "SuperDB2", Host: "MyDBServer", Port: 0, Name: "MyDB2", UsesIntegratedSecurity: false, ConnectionString: "connect to MyDB2"},
		},
		Loggers: []logger{
			logger{Name: "*", Level: "INFO", Target: "Gelf", Facility: "myservice", Destination: "logs.service.com"},
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

func TestParseMongoDBConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString string
		ExpectedCount    int
		ExpectedHost1    string
		ExpectedPort1    int
		ExpectedName     string
		ExpectedHost2    string
		ExpectedPort2    int
	}{
		{"mongodb://user:password@mongo1.company.com:2100,mongo2/db1?safe=true", 2, "mongo1.company.com", 2100, "db1", "mongo2", 0},
		{"mongodb://mongo1.company.com:2,mongo2/?safe=true", 2, "mongo1.company.com", 2, "", "mongo2", 0},
	} {
		dbs := parseMongoDBConnectionString(testCase.ConnectionString)

		if len(dbs) != testCase.ExpectedCount {
			t.Errorf("Unexpected db count, expected : %s actual : %s", testCase.ExpectedCount, len(dbs))
		}
		if dbs[0].Host != testCase.ExpectedHost1 {
			t.Errorf("Unexpected host1, expected : [%s] actual : [%s]", testCase.ExpectedHost1, dbs[0].Host)
		}
		if dbs[0].Port != testCase.ExpectedPort1 {
			t.Errorf("Unexpected port1, expected : %d actual : %d", testCase.ExpectedPort1, dbs[0].Port)
		}
		if dbs[0].Name != testCase.ExpectedName {
			t.Errorf("Unexpected name1, expected : %s actual : %s", testCase.ExpectedName, dbs[0].Name)
		}
		if dbs[0].UsesIntegratedSecurity {
			t.Errorf("Unexpected uses integrated security1, expected false but got true")
		}
		if dbs[0].ConnectionString != testCase.ConnectionString {
			t.Errorf("Unexpected connection string, expected : %s actual : %s", testCase.ConnectionString, dbs[0].ConnectionString)
		}
		if len(dbs) > 1 {
			if dbs[1].Host != testCase.ExpectedHost2 {
				t.Errorf("Unexpected host1, expected : [%s] actual : [%s]", testCase.ExpectedHost2, dbs[1].Host)
			}
			if dbs[1].Port != testCase.ExpectedPort2 {
				t.Errorf("Unexpected port1, expected : %d actual : %d", testCase.ExpectedPort2, dbs[1].Port)
			}
			if dbs[1].Name != testCase.ExpectedName {
				t.Errorf("Unexpected name1, expected : %s actual : %s", testCase.ExpectedName, dbs[1].Name)
			}
			if dbs[1].UsesIntegratedSecurity {
				t.Errorf("Unexpected uses integrated security1, expected false but got true")
			}
			if dbs[1].ConnectionString != testCase.ConnectionString {
				t.Errorf("Unexpected connection string, expected : %s actual : %s", testCase.ConnectionString, dbs[0].ConnectionString)
			}
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
			Loggers: []xmlNLogRulesLogger{
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
