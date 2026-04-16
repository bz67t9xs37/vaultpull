package config

import (
	"testing"
	"time"
)

func TestDefaultRetryConfig_Values(t *testing.T) {
	c := DefaultRetryConfig()
	if c.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", c.MaxAttempts)
	}
	if c.InitialDelay != 500*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %v", c.InitialDelay)
	}
	if c.MaxDelay != 10*time.Second {
		t.Errorf("unexpected MaxDelay: %v", c.MaxDelay)
	}
	if c.Multiplier != 2.0 {
		t.Errorf("unexpected Multiplier: %v", c.Multiplier)
	}
}

func TestApplyRetryDefaults_FillsZeroValues(t *testing.T) {
	c := &RetryConfig{}
	ApplyRetryDefaults(c)
	if c.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", c.MaxAttempts)
	}
	if c.InitialDelay != 500*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %v", c.InitialDelay)
	}
}

func TestApplyRetryDefaults_PreservesExistingValues(t *testing.T) {
	c := &RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   1.5,
	}
	ApplyRetryDefaults(c)
	if c.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts=5, got %d", c.MaxAttempts)
	}
	if c.Multiplier != 1.5 {
		t.Errorf("expected Multiplier=1.5, got %v", c.Multiplier)
	}
}

func TestApplyRetryDefaults_NilSafe(t *testing.T) {
	ApplyRetryDefaults(nil) // should not panic
}

func TestToPolicy_ReturnsFields(t *testing.T) {
	c := &RetryConfig{
		MaxAttempts:  4,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   3.0,
	}
	attempts, initial, max, mult := c.ToPolicy()
	if attempts != 4 || initial != 200*time.Millisecond || max != 5*time.Second || mult != 3.0 {
		t.Errorf("ToPolicy returned unexpected values: %v %v %v %v", attempts, initial, max, mult)
	}
}
