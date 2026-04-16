package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultNamespaceConfig_Values(t *testing.T) {
	cfg := config.DefaultNamespaceConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty prefix, got %q", cfg.Prefix)
	}
}

func TestApplyNamespaceDefaults_NilSafe(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyNamespaceDefaults(cfg)
	if cfg.Namespace == nil {
		t.Fatal("expected Namespace to be initialized")
	}
}

func TestApplyNamespaceDefaults_PreservesExistingPrefix(t *testing.T) {
	cfg := &config.Config{
		Namespace: &config.NamespaceConfig{
			Prefix: "prod",
		},
	}
	config.ApplyNamespaceDefaults(cfg)
	if cfg.Namespace.Prefix != "prod" {
		t.Errorf("expected prefix to be preserved, got %q", cfg.Namespace.Prefix)
	}
}

func TestApplyNamespaceDefaults_FillsStripMount(t *testing.T) {
	cfg := &config.Config{
		Namespace: &config.NamespaceConfig{},
	}
	config.ApplyNamespaceDefaults(cfg)
	if cfg.Namespace.StripMount != false {
		t.Errorf("expected StripMount default false")
	}
}
