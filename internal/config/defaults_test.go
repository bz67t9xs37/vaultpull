package config

import (
	"testing"
)

func TestApplyDefaults_Mount(t *testing.T) {
	cfg := &Config{
		Vault: VaultConfig{
			Address: "http://localhost:8200",
			Token:   "root",
		},
	}
	ApplyDefaults(cfg)
	if cfg.Vault.Mount != DefaultMount {
		t.Errorf("expected default mount %q, got %q", DefaultMount, cfg.Vault.Mount)
	}
}

func TestApplyDefaults_TokenFromEnv(t *testing.T) {
	original := lookupEnv
	defer func() { lookupEnv = original }()

	lookupEnv = func(key string) string {
		if key == "VAULT_TOKEN" {
			return "env-token"
		}
		return ""
	}

	cfg := &Config{
		Vault: VaultConfig{
			Address: "http://localhost:8200",
		},
	}
	ApplyDefaults(cfg)
	if cfg.Vault.Token != "env-token" {
		t.Errorf("expected token from env, got %q", cfg.Vault.Token)
	}
}

func TestApplyDefaults_AddressFromEnv(t *testing.T) {
	original := lookupEnv
	defer func() { lookupEnv = original }()

	lookupEnv = func(key string) string {
		if key == "VAULT_ADDR" {
			return "http://vault.example.com"
		}
		return ""
	}

	cfg := &Config{
		Vault: VaultConfig{
			Token: "root",
		},
	}
	ApplyDefaults(cfg)
	if cfg.Vault.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %q", cfg.Vault.Address)
	}
}

func TestApplyDefaults_ExistingValuesNotOverwritten(t *testing.T) {
	cfg := &Config{
		Vault: VaultConfig{
			Address: "http://custom:8200",
			Token:   "my-token",
			Mount:   "kv",
		},
	}
	ApplyDefaults(cfg)
	if cfg.Vault.Mount != "kv" {
		t.Errorf("mount should not be overwritten, got %q", cfg.Vault.Mount)
	}
	if cfg.Vault.Token != "my-token" {
		t.Errorf("token should not be overwritten, got %q", cfg.Vault.Token)
	}
}
