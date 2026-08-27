package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"nutrition-planner/config"
)

var openDB = sql.Open

// New builds a Database for cfg. dataDir is the directory sqlite database
// files are stored under (typically config.AppConfig.DataDirectory) - pass
// "" to keep the database file relative to the process's current working
// directory, or when cfg.Name is already a full path (e.g. in tests using
// t.TempDir()). dataDir is ignored for non-sqlite drivers, where cfg.Name
// is a database name inside the RDBMS, not a filesystem path.
func New(cfg config.DatabaseConfig, dataDir string) *Database {
	return &Database{
		cfg:     cfg,
		dataDir: dataDir,
	}
}

// Connect() establishes a connection to the database using the provided configuration. It builds the DSN (Data Source Name) based on the database driver and configuration, opens the database connection, and sets connection pooling parameters. If the connection is successful, it pings the database to ensure it's reachable. If any step fails, it returns an error.
func (db *Database) Connect() error {
	if db.usesDataDir() {
		if err := os.MkdirAll(db.dataDir, 0o755); err != nil {
			return fmt.Errorf("❌ Unable to create data directory %q: %w", db.dataDir, err)
		}
	}

	dsn, err := db.buildDSN()
	if err != nil {
		return err
	}

	conn, err := openDB(db.cfg.Driver, dsn)
	if err != nil {
		return fmt.Errorf("❌ Unable to open databse: %w", err)
	}

	// Set connection pooling parameters based on the database driver
	if db.cfg.Driver == "sqlite" || db.cfg.Driver == "sqlite3" {
		conn.SetMaxOpenConns(1)
		conn.SetMaxIdleConns(1)

		//Enable foreign key enforcement for SQLite
		if err := enableSQLiteFK(conn); err != nil {
			return err
		}
	} else {
		conn.SetMaxOpenConns(db.cfg.MaxOpenConns)
		conn.SetMaxIdleConns(db.cfg.MaxIdleConns)
	}
	conn.SetConnMaxLifetime(db.cfg.ConnMaxLifetime)
	conn.SetConnMaxIdleTime(db.cfg.ConnMaxIdleTime)

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("❌ Database is not reachable: %w", err)
	}

	db.conn = conn
	// Print a success message with the database connection details
	fmt.Printf(
		"✅ Connected to %s database %s@%s:%d\n",
		db.cfg.Driver,
		db.cfg.User,
		db.cfg.Host,
		db.cfg.Port,
	)

	return nil
}

// buildDSN constructs the Data Source Name (DSN) based on the database driver and configuration. It returns the DSN string or an error if the driver is unsupported.
func (db *Database) buildDSN() (string, error) {
	switch db.cfg.Driver {
	case "postgres", "pgx":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			db.cfg.Host,
			db.cfg.Port,
			db.cfg.User,
			db.cfg.Password,
			db.cfg.Name,
			db.cfg.SSLMode,
			db.cfg.Timezone,
		), nil

	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			db.cfg.User,
			db.cfg.Password,
			db.cfg.Host,
			db.cfg.Port,
			db.cfg.Name,
		), nil

	case "sqlite", "sqlite3":
		if db.cfg.Name == ":memory:" {
			return "file::memory:?cache=shared", nil
		}
		if db.usesDataDir() {
			return filepath.Join(db.dataDir, db.cfg.Name+".db"), nil
		}
		return fmt.Sprintf("%s.db", db.cfg.Name), nil

	default:
		return "", fmt.Errorf("❌ Unsupported database driver: %s", db.cfg.Driver)
	}
}

// All the query execution logic
func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	if db.conn == nil {
		return nil, fmt.Errorf("❌ Database connection is nil")
	}
	return db.conn.Query(query, args...)
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	if db.conn == nil {
		return nil, fmt.Errorf("❌ Database connection is nil")
	}
	return db.conn.Exec(query, args...)
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// Close the database connection
func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *Database) Driver() string {
	return db.cfg.Driver
}

// SQLitePath returns the on-disk path to the sqlite database file - used by
// settings.Handler.BackupDatabase to know what to copy. Returns an error
// for any non-sqlite driver (Postgres/MySQL have no single file to copy;
// backing those up is a different operation entirely) or for an in-memory
// database (nothing on disk to back up).
func (db *Database) SQLitePath() (string, error) {
	if db.cfg.Driver != "sqlite" && db.cfg.Driver != "sqlite3" {
		return "", fmt.Errorf("backup is only supported for sqlite databases (driver is %q)", db.cfg.Driver)
	}
	if db.cfg.Name == ":memory:" {
		return "", errors.New("cannot back up an in-memory database")
	}
	return db.buildDSN()
}

func enableSQLiteFK(conn *sql.DB) error {
	_, err := conn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("unable to enable SQLite foreign keys: %w", err)
	}

	return nil
}

// usesDataDir reports whether dataDir should be applied - only meaningful
// for a sqlite/sqlite3 file-backed database (not :memory:, and not for
// server databases like postgres/mysql where Name is a database name, not
// a filesystem path).
func (db *Database) usesDataDir() bool {
	return (db.cfg.Driver == "sqlite" || db.cfg.Driver == "sqlite3") &&
		db.cfg.Name != ":memory:" &&
		db.dataDir != ""
}
