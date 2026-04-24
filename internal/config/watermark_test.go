package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultWatermarkConfig_Values(t *testing.T) {
	cfg := DefaultWatermarkConfig()
	require.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.StoreDir, "default StoreDir should be set")
	assert.True(t, cfg.Enabled, "watermark should be enabled by default")
}

func TestApplyWatermarkDefaults_NilSafe(t *testing.T) {
	// Should not panic when called with a nil pointer receiver via a config
	// that has no watermark section set.
	cfg := &Config{}
	require.NotPanics(t, func() {
		ApplyWatermarkDefaults(cfg)
	})
	require.NotNil(t, cfg.Watermark)
}

func TestApplyWatermarkDefaults_FillsStoreDir(t *testing.T) {
	cfg := &Config{
		Watermark: &WatermarkConfig{},
	}
	ApplyWatermarkDefaults(cfg)
	assert.NotEmpty(t, cfg.Watermark.StoreDir)
}

func TestApplyWatermarkDefaults_PreservesExistingValues(t *testing.T) {
	customDir := "/custom/watermark"
	cfg := &Config{
		Watermark: &WatermarkConfig{
			StoreDir: customDir,
			Enabled:  true,
		},
	}
	ApplyWatermarkDefaults(cfg)
	assert.Equal(t, customDir, cfg.Watermark.StoreDir)
}

func TestIsEnabled_Watermark_True(t *testing.T) {
	cfg := &WatermarkConfig{Enabled: true}
	assert.True(t, cfg.IsEnabled())
}

func TestIsEnabled_Watermark_False(t *testing.T) {
	cfg := &WatermarkConfig{Enabled: false}
	assert.False(t, cfg.IsEnabled())
}
