package namespace_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/namespace"
)

func TestNamespace_UsesConfigPrefix(t *testing.T) {
	cfg := &config.Config{
		Namespace: &config.NamespaceConfig{
			Prefix: "secret/myapp",
			StripMount: false,
		},
	}

	r := namespace.New(cfg.Namespace.Prefix, cfg.Namespace.StripMount)
	resolved := r.Resolve("db/password")

	if resolved != "secret/myapp/db/password" {
		t.Errorf("unexpected resolved path: %s", resolved)
	}
}

func TestNamespace_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyNamespaceDefaults(cfg)

	r := namespace.New(cfg.Namespace.Prefix, cfg.Namespace.StripMount)
	resolved := r.Resolve("secret/foo")

	if resolved != "secret/foo" {
		t.Errorf("expected path unchanged without prefix, got %s", resolved)
	}
}

func TestNamespace_StripMountViaConfig(t *testing.T) {
	cfg := &config.Config{
		Namespace: &config.NamespaceConfig{
			Prefix:     "kv",
			StripMount: true,
		},
	}

	r := namespace.New(cfg.Namespace.Prefix, cfg.Namespace.StripMount)
	stripped := r.StripMount("kv/myapp/token")

	if stripped != "myapp/token" {
		t.Errorf("unexpected stripped path: %s", stripped)
	}
}
