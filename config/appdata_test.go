package config_test

import (
	"path/filepath"
	"testing"

	"nutrition-planner/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAppDataDir(t *testing.T) {
	dir, err := config.DefaultAppDataDir()
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.True(t, filepath.IsAbs(dir), "expected an absolute path, got %q", dir)
	assert.Equal(t, "NutriPlan", filepath.Base(dir))
}
