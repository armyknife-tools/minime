// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestDBConfigSimple(t *testing.T) {
	// Save original environment variables
	origDBType := os.Getenv("TOFU_REGISTRY_DB_TYPE")
	origDBURL := os.Getenv("TOFU_REGISTRY_DB_URL")

	// Restore environment variables after test
	defer func() {
		os.Setenv("TOFU_REGISTRY_DB_TYPE", origDBType)
		os.Setenv("TOFU_REGISTRY_DB_URL", origDBURL)
	}()

	t.Run("SQLite Default", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("TOFU_REGISTRY_DB_TYPE")
		os.Unsetenv("TOFU_REGISTRY_DB_URL")

		config := NewDBConfig(hclog.NewNullLogger())
		err := config.LoadFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "sqlite", config.Type)
	})

	t.Run("PostgreSQL Type", func(t *testing.T) {
		// Set PostgreSQL type
		os.Setenv("TOFU_REGISTRY_DB_TYPE", "postgres")
		os.Unsetenv("TOFU_REGISTRY_DB_URL")

		config := NewDBConfig(hclog.NewNullLogger())
		err := config.LoadFromEnv()
		assert.NoError(t, err)
		assert.Equal(t, "postgres", config.Type)
	})
}
