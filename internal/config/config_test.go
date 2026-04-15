package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
vault:
  address: http://127.0.0.1:8200
  token: root
  mount: secret
targets:
  - secret: myapp/prod
    env_file: .env
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address, got %q", cfg.Vault.Address)
	}
	if len(cfg.Targets) != 1 {
		t.Errorf("expected 1 target, got %d", len(cfg.Targets))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTemp(t, `
vault:
  token: root
targets:
  - secret: myapp/prod
    env_file: .env
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoad_NoTargets(t *testing.T) {
	path := writeTemp(t, `
vault:
  address: http://127.0.0.1:8200
  token: root
targets: []
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, `:::invalid yaml:::`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected YAML parse error")
	}
}
