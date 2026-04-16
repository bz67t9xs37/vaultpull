package config

import "testing"

func TestDefaultResolveConfig_Values(t *testing.T) {
	c := DefaultResolveConfig()
	if c.StripPrefix != "" {
		t.Errorf("expected empty StripPrefix, got %q", c.StripPrefix)
	}
	if c.AddPrefix != "" {
		t.Errorf("expected empty AddPrefix, got %q", c.AddPrefix)
	}
	if c.FlattenPath {
		t.Error("expected FlattenPath to be false")
	}
}

func TestApplyResolveDefaults_NilSafe(t *testing.T) {
	// Should not panic on nil input.
	ApplyResolveDefaults(nil)
}

func TestApplyResolveDefaults_PreservesExistingValues(t *testing.T) {
	c := &ResolveConfig{
		StripPrefix: "VAULT_",
		AddPrefix:   "APP_",
		FlattenPath: true,
	}
	ApplyResolveDefaults(c)
	if c.StripPrefix != "VAULT_" {
		t.Errorf("expected StripPrefix to be preserved, got %q", c.StripPrefix)
	}
	if c.AddPrefix != "APP_" {
		t.Errorf("expected AddPrefix to be preserved, got %q", c.AddPrefix)
	}
	if !c.FlattenPath {
		t.Error("expected FlattenPath to be preserved as true")
	}
}

func TestHasStripPrefix_True(t *testing.T) {
	c := &ResolveConfig{StripPrefix: "VAULT_"}
	if !c.HasStripPrefix() {
		t.Error("expected HasStripPrefix to return true")
	}
}

func TestHasStripPrefix_False(t *testing.T) {
	c := &ResolveConfig{}
	if c.HasStripPrefix() {
		t.Error("expected HasStripPrefix to return false")
	}
}

func TestHasAddPrefix_True(t *testing.T) {
	c := &ResolveConfig{AddPrefix: "APP_"}
	if !c.HasAddPrefix() {
		t.Error("expected HasAddPrefix to return true")
	}
}

func TestHasAddPrefix_NilSafe(t *testing.T) {
	var c *ResolveConfig
	if c.HasAddPrefix() {
		t.Error("expected HasAddPrefix on nil to return false")
	}
}
