package config_test

import (
	"testing"

	"nutrition-planner/config"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentIsolation_Production_NoSuffix(t *testing.T) {
	cfg := config.Config{
		App:      config.AppConfig{Environment: "production"},
		Database: config.DatabaseConfig{Name: "nutrition-planner"},
	}

	dbCfg, suffix := cfg.EnvironmentIsolation()
	assert.Equal(t, "nutrition-planner", dbCfg.Name)
	assert.Empty(t, suffix)
}

func TestEnvironmentIsolation_EmptyEnvironment_TreatedAsProduction(t *testing.T) {
	cfg := config.Config{
		App:      config.AppConfig{Environment: ""},
		Database: config.DatabaseConfig{Name: "nutrition-planner"},
	}

	dbCfg, suffix := cfg.EnvironmentIsolation()
	assert.Equal(t, "nutrition-planner", dbCfg.Name)
	assert.Empty(t, suffix)
}

func TestEnvironmentIsolation_NonProduction_AddsSuffix(t *testing.T) {
	tests := []string{"development", "testing", "staging"}
	for _, env := range tests {
		t.Run(env, func(t *testing.T) {
			cfg := config.Config{
				App:      config.AppConfig{Environment: env},
				Database: config.DatabaseConfig{Name: "nutrition-planner"},
			}

			dbCfg, suffix := cfg.EnvironmentIsolation()
			assert.Equal(t, "nutrition-planner-"+env, dbCfg.Name)
			assert.Equal(t, " ("+env+")", suffix)
		})
	}
}

func TestEnvironmentIsolation_DoesNotMutateOriginalConfig(t *testing.T) {
	cfg := config.Config{
		App:      config.AppConfig{Environment: "testing"},
		Database: config.DatabaseConfig{Name: "nutrition-planner"},
	}

	_, _ = cfg.EnvironmentIsolation()

	// The receiver is a value (Config, not *Config), so calling the method
	// must never leak a mutation back into the caller's cfg.Database.Name -
	// this pins that down explicitly, since it's exactly the kind of bug
	// that's easy to reintroduce if the method signature ever changes.
	assert.Equal(t, "nutrition-planner", cfg.Database.Name)
}

func TestEnvironmentIsolation_PreservesOtherDatabaseFields(t *testing.T) {
	cfg := config.Config{
		App: config.AppConfig{Environment: "testing"},
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			Name:   "nutrition-planner",
		},
	}

	dbCfg, _ := cfg.EnvironmentIsolation()
	assert.Equal(t, "sqlite", dbCfg.Driver)
}
