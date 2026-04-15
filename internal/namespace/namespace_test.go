package namespace_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/namespace"
)

func TestResolve_WithPrefix(t *testing.T) {
	r := namespace.New("myapp/prod", "secret")
	got := r.Resolve("database")
	want := "secret/data/myapp/prod/database"
	if got != want {
		t.Errorf("Resolve() = %q, want %q", got, want)
	}
}

func TestResolve_WithoutPrefix(t *testing.T) {
	r := namespace.New("", "secret")
	got := r.Resolve("database")
	want := "secret/data/database"
	if got != want {
		t.Errorf("Resolve() = %q, want %q", got, want)
	}
}

func TestResolve_TrimsSlashes(t *testing.T) {
	r := namespace.New("/myapp/", "/kv/")
	got := r.Resolve("/payments")
	want := "kv/data/myapp/payments"
	if got != want {
		t.Errorf("Resolve() = %q, want %q", got, want)
	}
}

func TestResolveAll_ReturnsMultiplePaths(t *testing.T) {
	r := namespace.New("app", "secret")
	names := []string{"db", "api", "cache"}
	got := r.ResolveAll(names)
	want := []string{
		"secret/data/app/db",
		"secret/data/app/api",
		"secret/data/app/cache",
	}
	if len(got) != len(want) {
		t.Fatalf("ResolveAll() returned %d paths, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ResolveAll()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStripMount_WithPrefix(t *testing.T) {
	r := namespace.New("myapp/prod", "secret")
	got := r.StripMount("secret/data/myapp/prod/database")
	want := "database"
	if got != want {
		t.Errorf("StripMount() = %q, want %q", got, want)
	}
}

func TestStripMount_WithoutPrefix(t *testing.T) {
	r := namespace.New("", "secret")
	got := r.StripMount("secret/data/database")
	want := "database"
	if got != want {
		t.Errorf("StripMount() = %q, want %q", got, want)
	}
}

func TestResolveAll_EmptyInput(t *testing.T) {
	r := namespace.New("app", "secret")
	got := r.ResolveAll([]string{})
	if len(got) != 0 {
		t.Errorf("ResolveAll([]) returned %d items, want 0", len(got))
	}
}
