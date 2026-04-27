package compare_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/compare"
)

func TestCompare_Equal(t *testing.T) {
	c := compare.New()
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "B": "2"}
	r := c.Compare("test", src, dst)
	if !r.Equal {
		t.Fatal("expected equal")
	}
}

func TestCompare_Added(t *testing.T) {
	c := compare.New()
	r := c.Compare("test", map[string]string{}, map[string]string{"NEW": "val"})
	if len(r.Added) != 1 || r.Added["NEW"] != "val" {
		t.Fatalf("expected NEW to be added, got %v", r.Added)
	}
	if r.Equal {
		t.Fatal("expected not equal")
	}
}

func TestCompare_Removed(t *testing.T) {
	c := compare.New()
	r := c.Compare("test", map[string]string{"OLD": "val"}, map[string]string{})
	if len(r.Removed) != 1 || r.Removed["OLD"] != "val" {
		t.Fatalf("expected OLD to be removed, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	c := compare.New()
	r := c.Compare("test",
		map[string]string{"K": "old"},
		map[string]string{"K": "new"},
	)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	cv := r.Changed["K"]
	if cv.Before != "old" || cv.After != "new" {
		t.Errorf("unexpected changed value: %+v", cv)
	}
}

func TestCompare_Summary_NoDiff(t *testing.T) {
	c := compare.New()
	r := c.Compare("secrets/app", map[string]string{"A": "1"}, map[string]string{"A": "1"})
	if !strings.Contains(r.Summary(), "no differences") {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_Summary_WithDiff(t *testing.T) {
	c := compare.New()
	r := c.Compare("secrets/app",
		map[string]string{"A": "1"},
		map[string]string{"B": "2"},
	)
	s := r.Summary()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "-1") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestCompare_Keys_Sorted(t *testing.T) {
	c := compare.New()
	r := c.Compare("p",
		map[string]string{"Z": "z", "A": "old"},
		map[string]string{"A": "new", "M": "m"},
	)
	keys := r.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("unexpected key order: %v", keys)
	}
}
