package main

import "testing"

func TestParseMsSqlConnectionString(t *testing.T) {
	for _, testCase := range []struct {
		ConnectionString               string
		ExpectedHost                   string
		ExpectedName                   string
		ExpectedUsesIntegratedSecurity bool
	}{
		{"Server=tcp:dbserver1; Database=myDB; MultiSubnetFailover=True; Integrated Security=SSPI;", "tcp:dbserver1", "myDB", true},
		{"Data Source=AAA; Initial Catalog=BBB; uid=auser; password=itsasecret", "AAA", "BBB", false},
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
