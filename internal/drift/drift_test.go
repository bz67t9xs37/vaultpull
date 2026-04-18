package drift_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/drift"
)

func TestCheck_NoDrift(t *testing.T) {
	d := drift.New(map[string]string{"FOO": "bar", "BAZ": "qux"})
	results := d.Check("/app/.env", map[string]string{"FOO": "bar", "BAZ": "qux"})
	for _, r := range results {
		if r.Drifted {
			t.Errorf("expected no drift for key %s", r.Key)
		}
	}
}

func TestCheck_ValueChanged(t *testing.T) {
	d := drift.New(map[string]string{"SECRET": "new"})
	results := d.Check("/app/.env", map[string]string{"SECRET": "old"})
	if len(results) != 1 || !results[0].Drifted {
		t.Fatal("expected drift for changed value")
	}
}

func TestCheck_MissingLocal(t *testing.T) {
	d := drift.New(map[string]string{"TOKEN": "abc"})
	results := d.Check("/app/.env", map[string]string{})
	if len(results) != 1 || !results[0].Drifted {
		t.Fatal("expected drift for missing local key")
	}
	if results[0].LocalValue != "" {
		t.Errorf("expected empty local value, got %q", results[0].LocalValue)
	}
}

func TestCheck_SetsPath(t *testing.T) {
	d := drift.New(map[string]string{"K": "v"})
	results := d.Check("/some/.env", map[string]string{"K": "v"})
	if results[0].Path != "/some/.env" {
		t.Errorf("unexpected path: %s", results[0].Path)
	}
}

func TestSummary_NoDrift(t *testing.T) {
	d := drift.New(map[string]string{"A": "1", "B": "2"})
	results := d.Check(".", map[string]string{"A": "1", "B": "2"})
	s := drift.Summary(results)
	if !strings.Contains(s, "no drift") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	d := drift.New(map[string]string{"A": "1", "B": "2"})
	results := d.Check(".", map[string]string{"A": "wrong", "B": "2"})
	s := drift.Summary(results)
	if !strings.Contains(s, "1/2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
