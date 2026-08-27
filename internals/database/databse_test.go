package database

import (
	"database/sql"
	"nutrition-planner/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func testConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Driver:       "postgres",
		Host:         "localhost",
		Port:         5432,
		User:         "test",
		Password:     "secret",
		Name:         "testdb",
		SSLMode:      "disable",
		Timezone:     "UTC",
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	}
}

func TestDB_Connect_Success(t *testing.T) {
	cfg := testConfig()

	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectPing()

	openDB = func(driver, dsn string) (*sql.DB, error) {
		return sqlDB, nil
	}
	defer func() { openDB = sql.Open }()

	db := New(cfg, "")

	err = db.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, db.conn)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_Connect_InvalidDriver(t *testing.T) {
	cfg := testConfig()
	cfg.Driver = "oracle"

	db := New(cfg, "")

	err := db.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "❌ Unsupported database driver")
}

func TestDB_Query(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	db := &Database{conn: sqlDB}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test")

	mock.ExpectQuery("SELECT id, name FROM users").
		WillReturnRows(rows)

	result, err := db.Query("SELECT id, name FROM users")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_Exec(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	db := &Database{conn: sqlDB}

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := db.Exec("INSERT INTO users(name) VALUES($1)", "John")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	rows, _ := result.RowsAffected()
	assert.Equal(t, int64(1), rows)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_Query_NoConnection(t *testing.T) {
	db := &Database{}

	rows, err := db.Query("SELECT 1")

	assert.Nil(t, rows)
	assert.Error(t, err)
}

func TestDB_Exec_NoConnection(t *testing.T) {
	db := &Database{}

	res, err := db.Exec("DELETE FROM users")

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestDB_Close(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	mock.ExpectClose()

	db := &Database{conn: sqlDB}

	err := db.Close()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestConnect_SQLiteForeignKeysEnabled(t *testing.T) {
	cfg := config.DatabaseConfig{
		Driver: "sqlite",
		Name:   ":memory:",
	}

	db := New(cfg, "")

	if err := db.Connect(); err != nil {
		t.Fatalf("Connect() failed: %v", err)
	}
	defer db.Close()

	var foreignKeys int

	if err := db.conn.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("failed to query PRAGMA foreign_keys: %v", err)
	}

	if foreignKeys != 1 {
		t.Fatalf("expected foreign_keys = 1, got %d", foreignKeys)
	}
}

func TestConnect_SQLiteForeignKeysAreEnforced(t *testing.T) {
	cfg := config.DatabaseConfig{
		Driver: "sqlite",
		Name:   ":memory:",
	}

	db := New(cfg, "")

	if err := db.Connect(); err != nil {
		t.Fatalf("Connect() failed: %v", err)
	}
	defer db.Close()

	_, err := db.Exec(`
		CREATE TABLE parents (
			id INTEGER PRIMARY KEY
		);

		CREATE TABLE children (
			id INTEGER PRIMARY KEY,
			parent_id INTEGER,
			FOREIGN KEY (parent_id) REFERENCES parents(id)
		);
	`)
	if err != nil {
		t.Fatalf("failed to create test tables: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO children (id, parent_id)
		VALUES (1, 999);
	`)

	if err == nil {
		t.Fatal("expected foreign key constraint error, got nil")
	}
}

func TestBuildDSN_SQLite_WithDataDir(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: "nutrition-planner"}
	db := New(cfg, "/some/data/dir")

	dsn, err := db.buildDSN()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/some/data/dir", "nutrition-planner.db"), dsn)
}

func TestBuildDSN_SQLite_WithoutDataDir(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: "nutrition-planner"}
	db := New(cfg, "")

	dsn, err := db.buildDSN()
	assert.NoError(t, err)
	assert.Equal(t, "nutrition-planner.db", dsn)
}

func TestBuildDSN_SQLite_InMemory_IgnoresDataDir(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: ":memory:"}
	db := New(cfg, "/some/data/dir")

	dsn, err := db.buildDSN()
	assert.NoError(t, err)
	assert.Equal(t, "file::memory:?cache=shared", dsn)
}

func TestBuildDSN_NonSQLite_IgnoresDataDir(t *testing.T) {
	cfg := testConfig() // driver: postgres
	db := New(cfg, "/some/data/dir")

	dsn, err := db.buildDSN()
	assert.NoError(t, err)
	// A data directory must never leak into a server DSN - Name there is a
	// database name inside the RDBMS, not a filesystem path.
	assert.NotContains(t, dsn, "/some/data/dir")
}

func TestConnect_SQLite_CreatesMissingDataDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "data")
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: "app"}
	db := New(cfg, dir)

	if err := db.Connect(); err != nil {
		t.Fatalf("Connect() failed: %v", err)
	}
	defer db.Close()

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected data directory to have been created: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", dir)
	}

	if _, err := os.Stat(filepath.Join(dir, "app.db")); err != nil {
		t.Fatalf("expected database file inside data directory: %v", err)
	}
}

func TestSQLitePath_ReturnsFileLocation(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: "app"}
	db := New(cfg, dir)

	path, err := db.SQLitePath()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "app.db"), path)
}

func TestSQLitePath_RejectsInMemory(t *testing.T) {
	cfg := config.DatabaseConfig{Driver: "sqlite", Name: ":memory:"}
	db := New(cfg, "")

	_, err := db.SQLitePath()
	assert.Error(t, err)
}

func TestSQLitePath_RejectsNonSQLiteDriver(t *testing.T) {
	db := New(testConfig(), "") // testConfig() uses driver "postgres"

	_, err := db.SQLitePath()
	assert.Error(t, err)
}
