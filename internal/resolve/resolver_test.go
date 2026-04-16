package resolve

import (
	"testing"
)

func TestResolve_SimpleKeys(t *testing.T) {
	r := New("secret")
	paths, err := r.Resolve([]string{"db/password", "api/key"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0].VaultPath != "secret/db/password" {
		t.Errorf("expected secret/db/password, got %s", paths[0].VaultPath)
	}
	if paths[1].VaultPath != "secret/api/key" {
		t.Errorf("expected secret/api/key, got %s", paths[1].VaultPath)
	}
}

func TestResolve_AlreadyQualifiedPath(t *testing.T) {
	r := New("secret")
	paths, err := r.Resolve([]string{"secret/db/password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if paths[0].VaultPath != "secret/db/password" {
		t.Errorf("path should not be duplicated: %s", paths[0].VaultPath)
	}
}

func TestResolve_LocalKeyName(t *testing.T) {
	r := New("secret")
	paths, err := r.Resolve([]string{"db/my-password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if paths[0].LocalKey != "MY_PASSWORD" {
		t.Errorf("expected MY_PASSWORD, got %s", paths[0].LocalKey)
	}
}

func TestResolve_EmptyKeys_ReturnsError(t *testing.T) {
	r := New("secret")
	_, err := r.Resolve([]string{})
	if err == nil {
		t.Fatal("expected error for empty keys, got nil")
	}
}

func TestResolve_AllBlankKeys_ReturnsError(t *testing.T) {
	r := New("secret")
	_, err := r.Resolve([]string{"   ", ""})
	if err == nil {
		t.Fatal("expected error when all keys are blank")
	}
}

func TestResolve_MountWithTrailingSlash(t *testing.T) {
	r := New("secret/")
	if r.Mount() != "secret" {
		t.Errorf("expected trailing slash to be trimmed, got %s", r.Mount())
	}
}

func TestLocalKeyName_SingleSegment(t *testing.T) {
	result := localKeyName("password")
	if result != "PASSWORD" {
		t.Errorf("expected PASSWORD, got %s", result)
	}
}

func TestLocalKeyName_MultiSegment(t *testing.T) {
	result := localKeyName("secret/database/host")
	if result != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %s", result)
	}
}
