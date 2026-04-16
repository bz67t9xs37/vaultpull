package config

import (
	"testing"
)

func TestDefaultMaskConfig_Values(t *testing.T) {
	c := DefaultMaskConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.MaskChar != "*" {
		t.Errorf("expected MaskChar '*', got %q", c.MaskChar)
	}
	if c.CustomKeys == nil {
		t.Error("expected CustomKeys to be non-nil")
	}
}

func TestApplyMaskDefaults_FillsMaskChar(t *testing.T) {
	c := &MaskConfig{Enabled: true}
	ApplyMaskDefaults(c)
	if c.MaskChar != "*" {
		t.Errorf("expected '*', got %q", c.MaskChar)
	}
}

func TestApplyMaskDefaults_FillsNilCustomKeys(t *testing.T) {
	c := &MaskConfig{MaskChar: "#"}
	ApplyMaskDefaults(c)
	if c.CustomKeys == nil {
		t.Error("expected CustomKeys to be initialised")
	}
}

func TestApplyMaskDefaults_PreservesExistingValues(t *testing.T) {
	c := &MaskConfig{
		MaskChar:   "X",
		CustomKeys: []string{"MY_SECRET"},
	}
	ApplyMaskDefaults(c)
	if c.MaskChar != "X" {
		t.Errorf("expected 'X', got %q", c.MaskChar)
	}
	if len(c.CustomKeys) != 1 || c.CustomKeys[0] != "MY_SECRET" {
		t.Error("custom keys were overwritten")
	}
}

func TestApplyMaskDefaults_NilSafe(t *testing.T) {
	// should not panic
	ApplyMaskDefaults(nil)
}

func TestHasCustomKeys_NoKeys(t *testing.T) {
	c := DefaultMaskConfig()
	if c.HasCustomKeys() {
		t.Error("expected HasCustomKeys to be false")
	}
}

func TestHasCustomKeys_WithKeys(t *testing.T) {
	c := &MaskConfig{CustomKeys: []string{"API_KEY"}}
	if !c.HasCustomKeys() {
		t.Error("expected HasCustomKeys to be true")
	}
}
