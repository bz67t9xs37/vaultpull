package sync_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/output"
	"github.com/yourusername/vaultpull/internal/sync"
)

type mockVault struct {
	secrets map[string]string
}

func (m *mockVault) GetSecrets(_ context.Context, _, _ string) (map[string]string, error) {
	return m.secrets, nil
}

func newTestConfig(t *testing.T, outputPath string, backup bool, bakDir string) *config.Config {
	t.Helper()
	return &config.Config{
		Address:   "http://127.0.0.1:8200",
		Token:     "test-token",
		Mount:     "secret",
		BackupDir: bakDir,
		Targets: []config.Target{
			{Path: "myapp/prod", Output: outputPath, Backup: backup},
		},
	}
}

func TestSync_NewSecrets(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")

	client := &mockVault{secrets: map[string]string{"DB_URL": "postgres://localhost"}}
	printer := output.New(os.Stdout)
	cfg := newTestConfig(t, outPath, false, "")

	s := sync.New(client, printer, cfg)
	if err := s.Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output file not written: %v", err)
	}
	if string(data) == "" {
		t.Error("expected non-empty .env file")
	}
}

func TestSync_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")
	_ = os.WriteFile(outPath, []byte("DB_URL=postgres://localhost\n"), 0o600)

	client := &mockVault{secrets: map[string]string{"DB_URL": "postgres://localhost"}}
	printer := output.New(os.Stdout)
	cfg := newTestConfig(t, outPath, false, "")

	s := sync.New(client, printer, cfg)
	if err := s.Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSync_BackupCreated(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, ".env")
	_ = os.WriteFile(outPath, []byte("OLD=1\n"), 0o600)
	bakDir := filepath.Join(tmpDir, "backups")

	client := &mockVault{secrets: map[string]string{"NEW": "2"}}
	printer := output.New(os.Stdout)
	cfg := newTestConfig(t, outPath, true, bakDir)

	s := sync.New(client, printer, cfg)
	if err := s.Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(bakDir)
	if err != nil {
		t.Fatalf("backup dir not created: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 backup file, got %d", len(entries))
	}
}
