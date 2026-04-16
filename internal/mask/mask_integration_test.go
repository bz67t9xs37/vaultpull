package mask_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/mask"
)

func TestMask_IntegratesWithConfig(t *testing.T) {
	cfg := config.DefaultMaskConfig()
	cfg.CustomKeys = []string{"CUSTOM_TOKEN"}

	m := mask.New(cfg.MaskChar, cfg.CustomKeys...)

	tests := []struct {
		key      string
		value    string
		expected string
	}{
		{"API_KEY", "super-secret", "************"},
		{"CUSTOM_TOKEN", "mytoken", "*******"},
		{"APP_NAME", "vaultpull", "vaultpull"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := m.MaskValue(tt.key, tt.value)
			if got != tt.expected {
				t.Errorf("key=%s: expected %q, got %q", tt.key, tt.expected, got)
			}
		})
	}
}

func TestMask_DisabledViaConfig(t *testing.T) {
	cfg := config.DefaultMaskConfig()
	cfg.Enabled = false

	if cfg.Enabled {
		t.Skip("masking is enabled, skipping disabled test")
	}

	// When masking is disabled the caller should skip masking entirely;
	// verify the config flag is readable.
	if cfg.MaskChar == "" {
		t.Error("MaskChar should still have a default even when disabled")
	}
}
