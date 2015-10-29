package main

import "testing"

func TestParseMsSqlConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString               string
		ExpectedSource                 string
		ExpectedDatabase               string
		ExpectedUsesIntegratedSecurity bool
	}{
		{"Server=tcp:dbserver1; Database=myDB; MultiSubnetFailover=True; Integrated Security=SSPI;", "tcp:dbserver1", "myDB", true},
		{"Data Source=AAA; Initial Catalog=BBB; uid=auser; password=itsasecret", "AAA", "BBB", false},
	} {
		db := parseMsSqlConnectionString(testCase.ConnectionString)

		if db.Source != testCase.ExpectedSource {
			t.Errorf("Unexpected source, expected : %s actual : %s", testCase.ExpectedSource, db.Source)
		}
		if db.Database != testCase.ExpectedDatabase {
			t.Errorf("Unexpected database, expected : %s actual : %s", testCase.ExpectedDatabase, db.Database)
		}
		if db.UsesIntegratedSecurity != testCase.ExpectedUsesIntegratedSecurity {
			t.Errorf("Unexpected uses integrated security: %t", testCase.ExpectedUsesIntegratedSecurity, db.UsesIntegratedSecurity)
		}
	}
}

func TestTransformNLogXml(t *testing.T) {
	source := xmlNLog{
		Targets: xmlNLogTargets{
			Targets: []xmlNLogTargetsTarget{
				xmlNLogTargetsTarget{
					Name:    "udpLogger",
					Type:    "NlogViewer",
					Address: "udp://localhost:7071",
				},
				xmlNLogTargetsTarget{
					Name:       "gelfLogger",
					Type:       "Gelf",
					Facility:   "super_service",
					GelfServer: "log.company.com",
					Port:       "12201",
				},
				xmlNLogTargetsTarget{
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
				xmlNLogRulesLogger{
					Name:     "*",
					MinLevel: "Trace",
					WriteTo:  "udpLogger",
				},
				xmlNLogRulesLogger{
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
