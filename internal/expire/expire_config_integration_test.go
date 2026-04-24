package expire_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/expire"
)

func TestExpire_UsesConfigWarningWindow(t *testing.T) {
	cfg := &config.ExpireConfig{
		Enabled:       true,
		WarningWindow: int64(48 * time.Hour),
	}

	checker := expire.New(time.Duration(cfg.WarningWindow))

	soon := time.Now().Add(24 * time.Hour)
	results := checker.Check(map[string]time.Time{
		"API_KEY": soon,
	})

	if len(results) == 0 {
		t.Fatal("expected at least one expiry result")
	}
	if !results[0].Warning {
		t.Error("expected Warning=true for key expiring within warning window")
	}
}

func TestExpire_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.Config{Expire: nil}
	config.ApplyExpireDefaults(cfg)

	if cfg.Expire == nil {
		t.Fatal("expected Expire config to be initialized")
	}

	checker := expire.New(time.Duration(cfg.Expire.WarningWindow))
	if checker == nil {
		t.Fatal("expected non-nil checker from default config")
	}
}

func TestExpire_DisabledViaConfig(t *testing.T) {
	cfg := &config.ExpireConfig{
		Enabled:       false,
		WarningWindow: int64(24 * time.Hour),
	}

	if cfg.IsEnabled() {
		t.Error("expected checker to be disabled via config")
	}

	// When disabled, no check should be run — simulate guard
	var results []expire.Result
	if cfg.IsEnabled() {
		checker := expire.New(time.Duration(cfg.WarningWindow))
		results = checker.Check(map[string]time.Time{
			"DB_PASS": time.Now().Add(-time.Hour),
		})
	}

	if len(results) != 0 {
		t.Errorf("expected no results when disabled, got %d", len(results))
	}
}
