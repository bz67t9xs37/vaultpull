package rotate_test

import (
	"testing"
	"time"

	"github.com/user/vaultpull/internal/rotate"
)

func TestIsDue_NeverRotated(t *testing.T) {
	s := &rotate.Schedule{Interval: 24 * time.Hour}
	if !s.IsDue() {
		t.Error("expected IsDue=true when LastRotated is zero")
	}
}

func TestIsDue_RecentRotation(t *testing.T) {
	s := &rotate.Schedule{
		Interval:    24 * time.Hour,
		LastRotated: time.Now().Add(-1 * time.Hour),
	}
	if s.IsDue() {
		t.Error("expected IsDue=false when rotation happened recently")
	}
}

func TestIsDue_OverdueRotation(t *testing.T) {
	s := &rotate.Schedule{
		Interval:    24 * time.Hour,
		LastRotated: time.Now().Add(-48 * time.Hour),
	}
	if !s.IsDue() {
		t.Error("expected IsDue=true when rotation is overdue")
	}
}

func TestNextRotation_Zero(t *testing.T) {
	s := &rotate.Schedule{Interval: time.Hour}
	next := s.NextRotation()
	if next.After(time.Now().Add(time.Second)) {
		t.Error("expected NextRotation to be approximately now when never rotated")
	}
}

func TestParseInterval_Hours(t *testing.T) {
	d, err := rotate.ParseInterval("6h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 6*time.Hour {
		t.Errorf("expected 6h, got %v", d)
	}
}

func TestParseInterval_Days(t *testing.T) {
	d, err := rotate.ParseInterval("7d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 7*24*time.Hour {
		t.Errorf("expected 168h, got %v", d)
	}
}

func TestParseInterval_Invalid(t *testing.T) {
	_, err := rotate.ParseInterval("banana")
	if err == nil {
		t.Error("expected error for invalid interval")
	}
}
