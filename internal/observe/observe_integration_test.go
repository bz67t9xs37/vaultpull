package observe_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/observe"
)

func newObserverFromConfig(c *config.ObserveConfig) *observe.Observer {
	if c == nil || !c.IsEnabled() {
		return nil
	}
	return observe.New()
}

func TestObserve_EnabledViaConfig(t *testing.T) {
	cfg := config.DefaultObserveConfig()
	cfg.Enabled = true

	o := newObserverFromConfig(&cfg)
	if o == nil {
		t.Fatal("expected observer to be created")
	}
	o.Info("secret/app", "sync started")
	if o.CountByLevel("INFO") != 1 {
		t.Errorf("expected 1 INFO event")
	}
}

func TestObserve_DisabledViaConfig(t *testing.T) {
	cfg := config.DefaultObserveConfig()
	cfg.Enabled = false

	o := newObserverFromConfig(&cfg)
	if o != nil {
		t.Error("expected nil observer when disabled")
	}
}

func TestObserve_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.ObserveConfig{Enabled: true}
	config.ApplyObserveDefaults(cfg)

	if cfg.Level != "INFO" {
		t.Errorf("expected default level INFO, got %s", cfg.Level)
	}
	if cfg.Target != "stderr" {
		t.Errorf("expected default target stderr, got %s", cfg.Target)
	}
}

func TestObserve_EmitsMultipleLevels(t *testing.T) {
	var buf bytes.Buffer
	o := observe.NewWithWriter(&buf)

	o.Info("secret/app", "synced")
	o.Warn("secret/db", "TTL low")
	o.Error("secret/cache", "fetch failed")

	out := buf.String()
	for _, want := range []string{"INFO", "WARN", "ERROR"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %s", want)
		}
	}
	if o.CountByLevel("INFO") != 1 || o.CountByLevel("WARN") != 1 || o.CountByLevel("ERROR") != 1 {
		t.Error("unexpected event counts")
	}
}
