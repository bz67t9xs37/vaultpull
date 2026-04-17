package audit_test

import (
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/config"
)

func TestAudit_UsesConfigLogPath(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.AuditConfig{
		Enabled:  true,
		LogPath:  filepath.Join(dir, "audit.log"),
		MaxLines: 100,
	}

	logger := audit.New(cfg.LogPath)
	err := logger.Record(audit.Entry{
		Path:    "secret/app",
		Action:  "sync",
		Changes: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error reading entries: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestAudit_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.AuditConfig{}
	config.ApplyAuditDefaults(cfg)

	if cfg.LogPath == "" {
		t.Error("expected LogPath to be set after applying defaults")
	}
	if cfg.MaxLines == 0 {
		t.Error("expected MaxLines to be set after applying defaults")
	}
}

func TestAudit_DisabledViaConfig(t *testing.T) {
	cfg := &config.AuditConfig{Enabled: false}
	if cfg.IsEnabled() {
		t.Error("expected audit to be disabled")
	}
}
