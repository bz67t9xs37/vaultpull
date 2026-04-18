package verify_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/verify"
)

func TestVerify_UsesConfigKeys(t *testing.T) {
	cfg := &config.VerifyConfig{
		Enabled:         true,
		Keys:            []string{"DB_URL", "API_KEY"},
		RequireNonEmpty: true,
	}
	config.ApplyVerifyDefaults(cfg)

	v := verify.New(cfg.RequireNonEmpty)
	secrets := map[string]string{
		"DB_URL":  "postgres://localhost",
		"API_KEY": "abc123",
	}

	results := v.Check(secrets, cfg.Keys)
	for _, r := range results {
		if r.Missing || r.Empty {
			t.Errorf("unexpected failure for key %s", r.Key)
		}
	}
}

func TestVerify_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.VerifyConfig{}
	config.ApplyVerifyDefaults(cfg)

	if cfg.Keys == nil {
		t.Error("expected Keys slice to be initialized after defaults")
	}

	v := verify.New(cfg.RequireNonEmpty)
	results := v.Check(map[string]string{"X": "y"}, cfg.Keys)
	if len(results) != 0 {
		t.Errorf("expected no results for empty key list, got %d", len(results))
	}
}

func TestVerify_MissingKeyDetectedViaConfig(t *testing.T) {
	cfg := &config.VerifyConfig{
		Enabled:         true,
		Keys:            []string{"MISSING_KEY"},
		RequireNonEmpty: false,
	}
	config.ApplyVerifyDefaults(cfg)

	v := verify.New(cfg.RequireNonEmpty)
	results := v.Check(map[string]string{}, cfg.Keys)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Missing {
		t.Error("expected key to be marked as missing")
	}
}
