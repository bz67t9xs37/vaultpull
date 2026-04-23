package expire_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

var now = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestCheck_NotExpired(t *testing.T) {
	c := expire.New(24 * time.Hour)
	c.(*expire.Checker) // ensure interface usage if needed
	_ = c

	checker := expire.New(24 * time.Hour)
	entries := []expire.Entry{
		{Path: "secret/app", ExpiresAt: now.Add(48 * time.Hour)},
	}
	results := checker.Check(entries)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Expired || results[0].Warning {
		t.Errorf("expected no expiry or warning for future secret")
	}
}

func TestCheck_Expired(t *testing.T) {
	checker := expire.New(24 * time.Hour)
	entries := []expire.Entry{
		{Path: "secret/old", ExpiresAt: now.Add(-1 * time.Hour)},
	}
	results := checker.Check(entries)
	if !results[0].Expired {
		t.Errorf("expected secret to be expired")
	}
}

func TestCheck_Warning(t *testing.T) {
	checker := expire.New(24 * time.Hour)
	entries := []expire.Entry{
		{Path: "secret/soon", ExpiresAt: now.Add(12 * time.Hour)},
	}
	results := checker.Check(entries)
	if results[0].Expired {
		t.Errorf("expected secret not yet expired")
	}
	if !results[0].Warning {
		t.Errorf("expected warning for soon-expiring secret")
	}
}

func TestCheck_ZeroExpirySkipped(t *testing.T) {
	checker := expire.New(24 * time.Hour)
	entries := []expire.Entry{
		{Path: "secret/no-ttl", ExpiresAt: time.Time{}},
	}
	results := checker.Check(entries)
	if len(results) != 0 {
		t.Errorf("expected zero-expiry entry to be skipped")
	}
}

func TestSummary_AllOk(t *testing.T) {
	results := []expire.Result{
		{Path: "a", Expired: false, Warning: false},
	}
	got := expire.Summary(results)
	if got != "all secrets are within their expiry window" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestSummary_Mixed(t *testing.T) {
	results := []expire.Result{
		{Path: "a", Expired: true},
		{Path: "b", Warning: true},
		{Path: "c"},
	}
	got := expire.Summary(results)
	if got != "1 expired, 1 expiring soon" {
		t.Errorf("unexpected summary: %s", got)
	}
}
