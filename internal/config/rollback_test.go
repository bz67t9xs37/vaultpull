package config_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/config"
)

func TestDefaultRollbackConfig_Values(t *testing.T) {
	c := config.DefaultRollbackConfig()
	if !c.Enabled {
		t.Error("expected Enabled=true")
	}
	if c.BackupDir == "" {
		t.Error("expected non-empty BackupDir")
	}
}

func TestApplyRollbackDefaults_NilSafe(t *testing.T) {
	// Must not panic
	config.ApplyRollbackDefaults(nil)
}

func TestApplyRollbackDefaults_FillsEmptyDir(t *testing.T) {
	c := &config.RollbackConfig{Enabled: true, BackupDir: ""}
	config.ApplyRollbackDefaults(c)
	if c.BackupDir == "" {
		t.Error("expected BackupDir to be filled")
	}
}

func TestApplyRollbackDefaults_PreservesExistingValues(t *testing.T) {
	c := &config.RollbackConfig{Enabled: true, BackupDir: "/custom/backups"}
	config.ApplyRollbackDefaults(c)
	if c.BackupDir != "/custom/backups" {
		t.Errorf("BackupDir overwritten, got %s", c.BackupDir)
	}
}

func TestIsEnabled_Rollback_True(t *testing.T) {
	c := &config.RollbackConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Rollback_False(t *testing.T) {
	c := &config.RollbackConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}

func TestIsEnabled_Rollback_Nil(t *testing.T) {
	var c *config.RollbackConfig
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false for nil config")
	}
}
