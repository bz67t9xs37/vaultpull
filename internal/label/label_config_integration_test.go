package label_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/label"
)

func TestLabel_UsesConfigStoreDir(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.LabelConfig{
		StoreDir: dir,
		Enabled:  true,
	}

	m := label.New(cfg.StoreDir)
	err := m.Set("secret/app", map[string]string{"env": "production"})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	labels, err := m.Get("secret/app")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if labels["env"] != "production" {
		t.Errorf("expected env=production, got %s", labels["env"])
	}
}

func TestLabel_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyLabelDefaults(cfg)

	if cfg.Label == nil {
		t.Fatal("expected Label config to be initialized")
	}
	if cfg.Label.StoreDir == "" {
		t.Error("expected StoreDir to be set after applying defaults")
	}
}

func TestLabel_DisabledViaConfig(t *testing.T) {
	cfg := &config.LabelConfig{
		Enabled: false,
	}

	if cfg.IsEnabled() {
		t.Error("expected label to be disabled")
	}
}

func TestLabel_MergesLabelsViaConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.LabelConfig{
		StoreDir: dir,
		Enabled:  true,
	}

	m := label.New(cfg.StoreDir)

	_ = m.Set("secret/app", map[string]string{"team": "platform"})
	_ = m.Set("secret/app", map[string]string{"env": "staging"})

	labels, err := m.Get("secret/app")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if labels["team"] != "platform" {
		t.Errorf("expected team=platform, got %s", labels["team"])
	}
	if labels["env"] != "staging" {
		t.Errorf("expected env=staging, got %s", labels["env"])
	}
}
