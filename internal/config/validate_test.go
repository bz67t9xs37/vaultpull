package config

import (
	"strings"
	"testing"
)

func validConfig() *Config {
	return &Config{
		Address: "http://127.0.0.1:8200",
		Token:   "s.testtoken",
		Mount:   "secret",
		Targets: []Target{
			{SecretPath: "myapp/prod", EnvFile: ".env"},
		},
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingAddress(t *testing.T) {
	cfg := validConfig()
	cfg.Address = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
	if !strings.Contains(err.Error(), "vault address is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_InvalidAddressURL(t *testing.T) {
	cfg := validConfig()
	cfg.Address = "not-a-url"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
	if !strings.Contains(err.Error(), "not a valid URL") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := validConfig()
	cfg.Token = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if !strings.Contains(err.Error(), "vault token is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MissingMount(t *testing.T) {
	cfg := validConfig()
	cfg.Mount = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing mount")
	}
	if !strings.Contains(err.Error(), "mount path is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_NoTargets(t *testing.T) {
	cfg := validConfig()
	cfg.Targets = nil
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for no targets")
	}
	if !strings.Contains(err.Error(), "at least one sync target") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_TargetMissingSecretPath(t *testing.T) {
	cfg := validConfig()
	cfg.Targets[0].SecretPath = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing secret_path")
	}
	if !strings.Contains(err.Error(), "secret_path is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &Config{}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected multiple errors")
	}
	if !IsValidationError(err) {
		t.Errorf("expected *ValidationError, got %T", err)
	}
	ve := err.(*ValidationError)
	if len(ve.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}
