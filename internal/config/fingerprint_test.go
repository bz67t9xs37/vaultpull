package config

import "testing"

func TestDefaultFingerprintConfig_Values(t *testing.T) {
	c := DefaultFingerprintConfig()
	if !c.Enabled {
		t.Error("expected Enabled=true by default")
	}
	if c.StoreDir == "" {
		t.Error("expected non-empty StoreDir")
	}
}

func TestApplyFingerprintDefaults_NilSafe(t *testing.T) {
	// should not panic
	ApplyFingerprintDefaults(nil)
}

func TestApplyFingerprintDefaults_FillsStoreDir(t *testing.T) {
	c := &FingerprintConfig{Enabled: true}
	ApplyFingerprintDefaults(c)
	if c.StoreDir == "" {
		t.Error("expected StoreDir to be filled")
	}
}

func TestApplyFingerprintDefaults_PreservesExistingValues(t *testing.T) {
	c := &FingerprintConfig{
		Enabled:  false,
		StoreDir: "/custom/fp",
	}
	ApplyFingerprintDefaults(c)
	if c.StoreDir != "/custom/fp" {
		t.Errorf("expected StoreDir preserved, got %s", c.StoreDir)
	}
}

func TestIsEnabled_Fingerprint_True(t *testing.T) {
	c := &FingerprintConfig{Enabled: true}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestIsEnabled_Fingerprint_False(t *testing.T) {
	c := &FingerprintConfig{Enabled: false}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false")
	}
}

func TestIsEnabled_Fingerprint_Nil(t *testing.T) {
	var c *FingerprintConfig
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false for nil config")
	}
}
