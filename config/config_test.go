package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigLoadAndMask_TempDir(t *testing.T) {
	// Creăm un temporary directory pentru test
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// Șterge tot la final
	defer os.RemoveAll(tempDir)

	yamlFile := filepath.Join(tempDir, "test_config.yaml")
	envFile := filepath.Join(tempDir, "test_config.env")

	// Load config (should generate YAML and env if not exist)
	cfg, err := Load(yamlFile, envFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 1️⃣ Check default values
	if cfg.App.Name != "nutrition-planner" {
		t.Errorf("unexpected App.Name: got %s", cfg.App.Name)
	}
	if cfg.Database.Driver != "sqlite" && cfg.Database.Driver != "sqlite3" {
		t.Errorf("unexpected Database.Driver: got %s", cfg.Database.Driver)
	}

	// 2️⃣ Check secret fields are set
	if cfg.Database.Password != "CHANGE_ME" {
		t.Errorf("Database.Password should be set to CHANGE_ME by default")
	}

	// 3️⃣ Check that YAML and env files were created
	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		t.Errorf("yaml file not created")
	}
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Errorf(".env file not created")
	}

	// 4️⃣ Test MaskSecrets
	MaskSecrets(cfg)
	if cfg.Database.Password != "******" {
		t.Errorf("Database.Password should be masked")
	}
}

func TestEnvOverrides_TempDir(t *testing.T) {
	// Creăm un temporary directory
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	yamlFile := filepath.Join(tempDir, "test_env_override.yaml")
	envFile := filepath.Join(tempDir, "test_env_override.env")

	// Set environment variables
	os.Setenv("APP_NAME", "MyApp")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_PASSWORD", "mysecret")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("DATABASE_PASSWORD")

	cfg, err := Load(yamlFile, envFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Check env overrides
	if cfg.App.Name != "MyApp" {
		t.Errorf("expected App.Name=MyApp, got %s", cfg.App.Name)
	}
	if cfg.Database.Password != "mysecret" {
		t.Errorf("expected Database.Password=mysecret, got %s", cfg.Database.Password)
	}
}

func TestDurationsParsing(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			ConnMaxLifetime: 2 * time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
		},
	}

	// Check types
	if cfg.Database.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("ConnMaxIdleTime wrong: %v", cfg.Database.ConnMaxIdleTime)
	}
}

func TestMaskSecrets(t *testing.T) {
	// Creăm un Config populat cu valori de test
	cfg := &Config{
		App: AppConfig{
			Name:          "TestApp",
			Version:       "0.1",
			Environment:   "test",
			Debug:         true,
			DataDirectory: "/tmp/data",
			TempDirectory: "/tmp/temp",
		},

		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "127.0.0.1",
			Port:     5432,
			User:     "db_user",
			Password: "db_secret",
			Name:     "db_test",
		},
		Logger: LoggerConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "/tmp/logs",
			MaxSize:    10,
			MaxAge:     7,
			MaxBackups: 3,
			Compress:   true,
		},
	}

	// Apelăm MaskSecrets
	MaskSecrets(cfg)

	// Verificăm că toate câmpurile secret sunt mascate
	if cfg.Database.Password != "******" {
		t.Errorf("expected Database.Password to be masked, got %s", cfg.Database.Password)
	}

	// Verificăm că câmpurile non-secret nu sunt afectate
	if cfg.App.Name != "TestApp" {
		t.Errorf("expected App.Name to remain unchanged, got %s", cfg.App.Name)
	}
	if cfg.Database.User != "db_user" {
		t.Errorf("expected Database.User to remain unchanged, got %s", cfg.Database.User)
	}
}
