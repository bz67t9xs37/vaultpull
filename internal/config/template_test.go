package config

import "testing"

func TestDefaultTemplateConfig_Values(t *testing.T) {
	c := DefaultTemplateConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Dir != "templates" {
		t.Errorf("unexpected Dir: %s", c.Dir)
	}
	if c.LeftDelim != "{{" {
		t.Errorf("unexpected LeftDelim: %s", c.LeftDelim)
	}
	if c.RightDelim != "}}" {
		t.Errorf("unexpected RightDelim: %s", c.RightDelim)
	}
}

func TestApplyTemplateDefaults_FillsEmptyDir(t *testing.T) {
	c := &TemplateConfig{}
	ApplyTemplateDefaults(c)
	if c.Dir != "templates" {
		t.Errorf("expected Dir to be filled, got %s", c.Dir)
	}
}

func TestApplyTemplateDefaults_FillsDelimiters(t *testing.T) {
	c := &TemplateConfig{}
	ApplyTemplateDefaults(c)
	if c.LeftDelim != "{{" {
		t.Errorf("expected LeftDelim to be filled, got %s", c.LeftDelim)
	}
	if c.RightDelim != "}}" {
		t.Errorf("expected RightDelim to be filled, got %s", c.RightDelim)
	}
}

func TestApplyTemplateDefaults_PreservesExistingValues(t *testing.T) {
	c := &TemplateConfig{
		Dir:        "custom-templates",
		LeftDelim:  "[[" ,
		RightDelim: "]]",
	}
	ApplyTemplateDefaults(c)
	if c.Dir != "custom-templates" {
		t.Errorf("Dir overwritten: %s", c.Dir)
	}
	if c.LeftDelim != "[[" {
		t.Errorf("LeftDelim overwritten: %s", c.LeftDelim)
	}
}

func TestApplyTemplateDefaults_NilSafe(t *testing.T) {
	ApplyTemplateDefaults(nil) // should not panic
}

func TestIsEnabled_True(t *testing.T) {
	c := &TemplateConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_False(t *testing.T) {
	c := &TemplateConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}

func TestIsEnabled_Nil(t *testing.T) {
	var c *TemplateConfig
	if c.IsEnabled() {
		t.Error("expected nil TemplateConfig to return false for IsEnabled")
	}
}
