package transform_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/transform"
)

func TestApply_Prefix(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "prefix", To: "prod_"},
	})
	out, err := tr.Apply(map[string]string{"DB_PASS": "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "prod_secret" {
		t.Errorf("expected prod_secret, got %s", out["DB_PASS"])
	}
}

func TestApply_Suffix(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "suffix", To: "_v2"},
	})
	out, err := tr.Apply(map[string]string{"API_KEY": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "abc123_v2" {
		t.Errorf("expected abc123_v2, got %s", out["API_KEY"])
	}
}

func TestApply_Uppercase(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "uppercase"},
	})
	out, err := tr.Apply(map[string]string{"HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "LOCALHOST" {
		t.Errorf("expected LOCALHOST, got %s", out["HOST"])
	}
}

func TestApply_Replace(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "replace", From: "staging", To: "production"},
	})
	out, err := tr.Apply(map[string]string{"ENV": "staging", "URL": "https://staging.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENV"] != "production" {
		t.Errorf("expected production, got %s", out["ENV"])
	}
	if out["URL"] != "https://production.example.com" {
		t.Errorf("expected https://production.example.com, got %s", out["URL"])
	}
}

func TestApply_ScopedToKey(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "uppercase", Key: "DB_PASS"},
	})
	out, err := tr.Apply(map[string]string{"DB_PASS": "secret", "API_KEY": "mykey"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "SECRET" {
		t.Errorf("expected SECRET, got %s", out["DB_PASS"])
	}
	if out["API_KEY"] != "mykey" {
		t.Errorf("expected mykey unchanged, got %s", out["API_KEY"])
	}
}

func TestApply_UnknownRuleType_ReturnsError(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "rot13"},
	})
	_, err := tr.Apply(map[string]string{"KEY": "value"})
	if err == nil {
		t.Fatal("expected error for unknown rule type, got nil")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{Type: "uppercase"},
	})
	out, err := tr.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
