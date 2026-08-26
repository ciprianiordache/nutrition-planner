package database

import (
	"database/sql"
	"nutrition-planner/config"
)

// Database represents the database connection and configuration.
type Database struct {
	cfg     config.DatabaseConfig // database configuration - loaded from the config package
	dataDir string                // directory sqlite database files are stored under (config.AppConfig.DataDirectory); empty means "current working directory", matching the old behavior
	conn    *sql.DB
}
