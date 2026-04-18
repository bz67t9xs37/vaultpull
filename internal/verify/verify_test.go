package verify_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/verify"
)

func TestCheck_AllPresent(t *testing.T) {
	v := verify.New(true)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	report, err := v.Check("secret/app", secrets, []string{"FOO", "BAZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasErrors() {
		t.Error("expected no errors")
	}
	if len(report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(report.Results))
	}
}

func TestCheck_MissingKey(t *testing.T) {
	v := verify.New(false)
	secrets := map[string]string{"FOO": "bar"}
	report, err := v.Check("secret/app", secrets, []string{"FOO", "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasErrors() {
		t.Error("expected errors due to missing key")
	}
}

func TestCheck_EmptyValue_RequireNonEmpty(t *testing.T) {
	v := verify.New(true)
	secrets := map[string]string{"FOO": ""}
	report, err := v.Check("secret/app", secrets, []string{"FOO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 1 || !report.Results[0].Empty {
		t.Error("expected FOO to be flagged as empty")
	}
}

func TestCheck_EmptyValue_NotRequired(t *testing.T) {
	v := verify.New(false)
	secrets := map[string]string{"FOO": ""}
	report, err := v.Check("secret/app", secrets, []string{"FOO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Results[0].Empty {
		t.Error("expected FOO not to be flagged as empty when requireNonEmpty=false")
	}
}

func TestCheck_NoExpectedKeys_ReturnsError(t *testing.T) {
	v := verify.New(true)
	_, err := v.Check("secret/app", map[string]string{}, []string{})
	if err == nil {
		t.Error("expected error for empty expectedKeys")
	}
}

func TestSummary_AllOk(t *testing.T) {
	v := verify.New(true)
	secrets := map[string]string{"A": "1", "B": "2"}
	report, _ := v.Check("secret/svc", secrets, []string{"A", "B"})
	summary := report.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_WithMissing(t *testing.T) {
	v := verify.New(false)
	secrets := map[string]string{}
	report, _ := v.Check("secret/svc", secrets, []string{"X"})
	summary := report.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
