package diff

import (
	"testing"
)

func TestCompute_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"FOO": "bar"}

	changes := Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Added {
		t.Errorf("expected Added, got %s", changes[0].Type)
	}
}

func TestCompute_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{}

	changes := Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", changes[0].Type)
	}
}

func TestCompute_Modified(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new := map[string]string{"FOO": "new"}

	changes := Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", changes[0].Type)
	}
	if changes[0].OldVal != "old" || changes[0].NewVal != "new" {
		t.Errorf("unexpected values: old=%s new=%s", changes[0].OldVal, changes[0].NewVal)
	}
}

func TestCompute_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar"}

	changes := Compute(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Unchanged {
		t.Errorf("expected Unchanged, got %s", changes[0].Type)
	}
}

func TestSummary(t *testing.T) {
	changes := []Change{
		{Type: Added},
		{Type: Added},
		{Type: Removed},
		{Type: Modified},
		{Type: Unchanged},
	}
	got := Summary(changes)
	expected := "+2 added, -1 removed, ~1 modified"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
