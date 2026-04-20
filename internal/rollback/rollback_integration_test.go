package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/rollback"
)

func TestRollback_UsesConfigBackupDir(t *testing.T) {
	dir := tempDir(t)
	cfg := &config.RollbackConfig{
		Enabled:   true,
		BackupDir: dir,
	}
	config.ApplyRollbackDefaults(cfg)

	// Write a backup file into the configured dir.
	const want = "API_KEY=supersecret\n"
	writeBackup(t, cfg.BackupDir, ".env.2024-06-01T10:00:00.bak", want)

	target := filepath.Join(dir, ".env")
	_ = os.WriteFile(target, []byte("API_KEY=stale\n"), 0o600)

	rb := rollback.New(cfg.BackupDir)
	entry, err := rb.Restore(target)
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	got, _ := os.ReadFile(target)
	if string(got) != want {
		t.Errorf("want %q, got %q", want, string(got))
	}
}

func TestRollback_ConfigDefaultsApplied(t *testing.T) {
	cfg := &config.RollbackConfig{}
	config.ApplyRollbackDefaults(cfg)
	if cfg.BackupDir == "" {
		t.Error("expected BackupDir to be set by defaults")
	}
}

func TestRollback_DisabledViaConfig(t *testing.T) {
	cfg := &config.RollbackConfig{Enabled: false, BackupDir: ".vaultpull/backups"}
	if cfg.IsEnabled() {
		t.Error("rollback should be disabled")
	}
	// When disabled, callers should skip rollback — nothing to assert on
	// the rollbacker itself, but we verify the flag is respected.
}
