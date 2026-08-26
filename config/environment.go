package config

import "fmt"

// EnvironmentIsolation derives the effective DatabaseConfig and an export
// folder suffix from c.App.Environment, so that anything other than
// "production" (an empty Environment counts as "production" too) gets its
// own database file and export folder - keeping testing/experimenting
// from ever touching or mixing with real client data.
//
// Every entry point that opens the database or resolves the export
// directory (the Wails app in app.go, cmd/seed, any future cmd/ tool)
// MUST go through this method instead of using c.Database / a raw export
// path directly - that's what keeps them all in agreement about which
// database file "the testing environment" actually means. Before this
// existed, app.go alone applied the suffix while cmd/seed didn't, so
// seeded data silently landed in a different .db file than the one the
// running app read from.
//
// dbCfg is c.Database with Name suffixed (e.g. "nutrition-planner" becomes
// "nutrition-planner-testing"). exportSuffix is a string to append to the
// export base directory (e.g. " (testing)"), or "" for production.
func (c Config) EnvironmentIsolation() (dbCfg DatabaseConfig, exportSuffix string) {
	dbCfg = c.Database
	if c.App.Environment != "" && c.App.Environment != "production" {
		exportSuffix = fmt.Sprintf(" (%s)", c.App.Environment)
		dbCfg.Name = fmt.Sprintf("%s-%s", dbCfg.Name, c.App.Environment)
	}
	return dbCfg, exportSuffix
}
