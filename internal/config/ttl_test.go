package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vaultpull/internal/config"
)

func TestDefaultTTLConfig_Values(t *testing.T) {
	cfg := config.DefaultTTLConfig()
	require.NotNil(t, cfg)
	assert.Greater(t, cfg.Default time.Duration(0))
}

func TestApplyTTLDefaults_NilSafe(t *testing.T) {
	cfg := &config.Config{
		Address: "http://vault:8200",
		Token:   "tok",
		Targets: "secret/app", Output: ".env"}},
		TTL:     nil,
	}
	config.ApplyTTLDefaults(cfg)
	assert.NotNil(t, cfg.TTL)
}

func TestApplyTTLDefaults_FillsZeroDuration(t *testing.T) {
	cfg := &config.Config{
		TTL: &config.TTLConfig{Default: 0},
	}
	config.ApplyTTLDefaults(cfg)
	assert.Greater(t, cfg.TTL.Default, time.Duration(0))
}

func TestApplyTTLDefaults_PreservesExistingValues(t *testing.T) {
	custom := 99 * time.Minute
	cfg := &config.Config{
		TTL: &config.TTLConfig{Default: custom},
	}
	config.ApplyTTLDefaults(cfg)
	assert.Equal(t, custom, cfg.TTL.Default)
}

func TestTTLConfig_IntegratesWithApplyDefaults(t *testing.T) {
	cfg := &config.Config{
		Address: "http://vault:8200",
		Token:   "tok",
		Targets: []config.Target{{Path: "secret/app", Output: ".env"}},
	}
	config.ApplyDefaults(cfg)
	assert.NotNil(t, cfg.TTL)
	assert.Greater(t, cfg.TTL.Default, time.Duration(0))
}
