package mask_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/mask"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	m := mask.New(nil, "")

	sensitive := []string{"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "PRIVATE_KEY", "aws_api_key"}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NonSensitiveKeys(t *testing.T) {
	m := mask.New(nil, "")

	public := []string{"APP_NAME", "LOG_LEVEL", "PORT", "HOST", "REGION"}
	for _, key := range public {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMaskValue_RedactsSensitive(t *testing.T) {
	m := mask.New(nil, "")

	got := m.MaskValue("DB_PASSWORD", "supersecret")
	if got == "supersecret" {
		t.Error("expected value to be masked, got original value")
	}
	if got == "" {
		t.Error("expected non-empty mask string")
	}
}

func TestMaskValue_LeavesPublicUnchanged(t *testing.T) {
	m := mask.New(nil, "")

	got := m.MaskValue("APP_NAME", "vaultpull")
	if got != "vaultpull" {
		t.Errorf("expected %q, got %q", "vaultpull", got)
	}
}

func TestMaskValue_CustomMaskChar(t *testing.T) {
	m := mask.New(nil, "[REDACTED]")

	got := m.MaskValue("API_SECRET", "abc123")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMaskMap_RedactsOnlySensitiveKeys(t *testing.T) {
	m := mask.New(nil, "")

	input := map[string]string{
		"APP_NAME":    "vaultpull",
		"DB_PASSWORD": "s3cr3t",
		"PORT":        "8080",
		"AUTH_TOKEN":  "tok_xyz",
	}

	result := m.MaskMap(input)

	if result["APP_NAME"] != "vaultpull" {
		t.Errorf("APP_NAME should be unchanged, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] == "s3cr3t" {
		t.Error("DB_PASSWORD should be masked")
	}
	if result["AUTH_TOKEN"] == "tok_xyz" {
		t.Error("AUTH_TOKEN should be masked")
	}
}

func TestMaskMap_CustomSensitiveKeys(t *testing.T) {
	m := mask.New([]string{"internal"}, "***")

	input := map[string]string{
		"INTERNAL_URL": "http://internal.svc",
		"PUBLIC_URL":   "http://example.com",
	}

	result := m.MaskMap(input)

	if result["INTERNAL_URL"] != "***" {
		t.Errorf("expected INTERNAL_URL to be masked, got %q", result["INTERNAL_URL"])
	}
	if result["PUBLIC_URL"] != "http://example.com" {
		t.Errorf("expected PUBLIC_URL unchanged, got %q", result["PUBLIC_URL"])
	}
}
