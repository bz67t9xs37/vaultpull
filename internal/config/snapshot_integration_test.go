package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vaultpull/internal/config"
)

func TestSnapshotConfig_IntegratesWithApplyDefaults(t *testing.T) {
	cfg := &config.Config{
		Address: "http://vault:8200",
		Token:   "tok",
		Targets: []config.Target{{Path: "secret/app", Output: ".env"}},
	}

	config.ApplyDefaults(cfg)

	assert.NotNil(t, cfg.Snapshot)
	assert.NotEmpty(t, cfg.Snapshot.Dir)
	assert.Greater(t, cfg.Snapshot.MaxPerPath, 0)
}

func TestSnapshotConfig_PreservesUserValues(t *testing.T) {
	cfg := &config.Config{
		Address: "http://vault:8200",
		Token:   "tok",
		Targets: []config.Target{{Path: "secret/app", Output: ".env"}},
		Snapshot: &config.SnapshotConfig{
			Dir:        "/custom/snapshots",
			MaxPerPath: 10,
		},
	}

	config.ApplyDefaults(cfg)

	assert.Equal(t, "/custom/snapshots", cfg.Snapshot.Dir)
	assert.Equal(t, 10, cfg.Snapshot.MaxPerPath)
}
