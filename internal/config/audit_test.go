package config

import (
	"testing"
)

func TestDefaultAuditConfig_Values(t *testing.T) {
	c := DefaultAuditConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.LogPath == "" {
		t.Error("expected non-empty LogPath")
	}
	if c.MaxLines <= 0 {
		t.Errorf("expected MaxLines > 0, got %d", c.MaxLines)
	}
}

func TestApplyAuditDefaults_FillsLogPath(t *testing.T) {
	c := &AuditConfig{}
	ApplyAuditDefaults(c)
	if c.LogPath == "" {
		t.Error("expected LogPath to be filled")
	}
}

func TestApplyAuditDefaults_FillsMaxLines(t *testing.T) {
	c := &AuditConfig{}
	ApplyAuditDefaults(c)
	if c.MaxLines == 0 {
		t.Error("expected MaxLines to be filled")
	}
}

func TestApplyAuditDefaults_PreservesExistingValues(t *testing.T) {
	c := &AuditConfig{
		LogPath:  "custom.log",
		MaxLines: 42,
	}
	ApplyAuditDefaults(c)
	if c.LogPath != "custom.log" {
		t.Errorf("expected LogPath to be preserved, got %s", c.LogPath)
	}
	if c.MaxLines != 42 {
		t.Errorf("expected MaxLines to be preserved, got %d", c.MaxLines)
	}
}

func TestApplyAuditDefaults_NilSafe(t *testing.T) {
	ApplyAuditDefaults(nil) // should not panic
}

func TestIsEnabled_True(t *testing.T) {
	c := &AuditConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_False(t *testing.T) {
	c := &AuditConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}

func TestIsEnabled_Nil(t *testing.T) {
	var c *AuditConfig
	if c.IsEnabled() {
		t.Error("expected IsEnabled to return false for nil config")
	}
}
