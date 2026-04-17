package config

import (
	"testing"
)

func TestDefaultOutputConfig_Values(t *testing.T) {
	c := DefaultOutputConfig()
	if c.Format != OutputFormatText {
		t.Errorf("expected format %q, got %q", OutputFormatText, c.Format)
	}
	if c.Verbose {
		t.Error("expected verbose to be false")
	}
	if c.NoColor {
		t.Error("expected no_color to be false")
	}
	if !c.MaskSensitive {
		t.Error("expected mask_sensitive to be true")
	}
}

func TestApplyOutputDefaults_FillsFormat(t *testing.T) {
	c := &OutputConfig{}
	ApplyOutputDefaults(c)
	if c.Format != OutputFormatText {
		t.Errorf("expected %q, got %q", OutputFormatText, c.Format)
	}
}

func TestApplyOutputDefaults_PreservesExistingValues(t *testing.T) {
	c := &OutputConfig{Format: OutputFormatJSON, Verbose: true}
	ApplyOutputDefaults(c)
	if c.Format != OutputFormatJSON {
		t.Errorf("expected format to be preserved as %q", OutputFormatJSON)
	}
	if !c.Verbose {
		t.Error("expected verbose to remain true")
	}
}

func TestApplyOutputDefaults_NilSafe(t *testing.T) {
	ApplyOutputDefaults(nil) // should not panic
}

func TestIsJSON_True(t *testing.T) {
	c := &OutputConfig{Format: OutputFormatJSON}
	if !c.IsJSON() {
		t.Error("expected IsJSON to return true")
	}
}

func TestIsJSON_False(t *testing.T) {
	c := &OutputConfig{Format: OutputFormatText}
	if c.IsJSON() {
		t.Error("expected IsJSON to return false")
	}
}

func TestIsVerbose_True(t *testing.T) {
	c := &OutputConfig{Verbose: true}
	if !c.IsVerbose() {
		t.Error("expected IsVerbose to return true")
	}
}

func TestIsVerbose_NilSafe(t *testing.T) {
	var c *OutputConfig
	if c.IsVerbose() {
		t.Error("expected nil OutputConfig to return false for IsVerbose")
	}
}
