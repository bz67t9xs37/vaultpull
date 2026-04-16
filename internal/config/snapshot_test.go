package config_test

import (
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultSnapshotConfig_Values(t *testing.T) {
	d := config.DefaultSnapshotConfig()
	if !d.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	want := filepath.Join(".vaultpull", "snapshots")
	if d.Dir != want {
		t.Errorf("expected Dir %q, got %q", want, d.Dir)
	}
	if d.MaxPerPath != 10 {
		t.Errorf("expected MaxPerPath 10, got %d", d.MaxPerPath)
	}
}

func TestApplySnapshotDefaults_FillsEmptyDir(t *testing.T) {
	c := config.SnapshotConfig{Enabled: true}
	config.ApplySnapshotDefaults(&c)
	want := filepath.Join(".vaultpull", "snapshots")
	if c.Dir != want {
		t.Errorf("expected Dir %q, got %q", want, c.Dir)
	}
}

func TestApplySnapshotDefaults_FillsZeroMaxPerPath(t *testing.T) {
	c := config.SnapshotConfig{Dir: "/tmp/snaps"}
	config.ApplySnapshotDefaults(&c)
	if c.MaxPerPath != 10 {
		t.Errorf("expected MaxPerPath 10, got %d", c.MaxPerPath)
	}
}

func TestApplySnapshotDefaults_PreservesExistingValues(t *testing.T) {
	c := config.SnapshotConfig{
		Enabled:    false,
		Dir:        "/custom/path",
		MaxPerPath: 5,
	}
	config.ApplySnapshotDefaults(&c)
	if c.Dir != "/custom/path" {
		t.Errorf("Dir should not be overwritten, got %q", c.Dir)
	}
	if c.MaxPerPath != 5 {
		t.Errorf("MaxPerPath should not be overwritten, got %d", c.MaxPerPath)
	}
	if c.Enabled {
		t.Error("Enabled should remain false")
	}
}
