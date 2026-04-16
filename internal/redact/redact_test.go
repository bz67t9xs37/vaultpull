package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/redact"
)

func TestRedact_LiteralSecret(t *testing.T) {
	r := redact.New([]string{"s3cr3t", "my-token"})

	got := r.Redact("password is s3cr3t and token is my-token")
	want := "password is [REDACTED] and token is [REDACTED]"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_ShortValueIgnored(t *testing.T) {
	r := redact.New([]string{"ab", "abc"})

	got := r.Redact("abc ab")
	if got != "abc ab" {
		t.Errorf("expected no redaction for short values, got %q", got)
	}
}

func TestRedact_PatternRule(t *testing.T) {
	r := redact.New(nil)
	err := r.AddPattern(`hvs\.[A-Za-z0-9]+`, "[VAULT_TOKEN]")
	if err != nil {
		t.Fatalf("AddPattern error: %v", err)
	}

	got := r.Redact("token=hvs.ABCdef123")
	want := "token=[VAULT_TOKEN]"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_InvalidPattern(t *testing.T) {
	r := redact.New(nil)
	err := r.AddPattern(`[invalid`, "x")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestRedactMap_RedactsValues(t *testing.T) {
	r := redact.New([]string{"topsecret"})

	input := map[string]string{
		"API_KEY":  "topsecret",
		"APP_NAME": "vaultpull",
	}

	out := r.RedactMap(input)

	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY: got %q, want [REDACTED]", out["API_KEY"])
	}
	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("APP_NAME: got %q, want %q", out["APP_NAME"], "vaultpull")
	}
}

func TestRedactMap_DoesNotMutateOriginal(t *testing.T) {
	r := redact.New([]string{"secret123"})

	input := map[string]string{"KEY": "secret123"}
	_ = r.RedactMap(input)

	if input["KEY"] != "secret123" {
		t.Error("original map was mutated")
	}
}
