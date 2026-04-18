package promote_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/promote"
)

type mockStore struct {
	data     map[string]map[string]string
	written  map[string]map[string]string
	fetchErr error
	writeErr error
}

func (m *mockStore) GetSecrets(_ context.Context, path string) (map[string]string, error) {
	if m.fetchErr != nil {
		return nil, m.fetchErr
	}
	return m.data[path], nil
}

func (m *mockStore) WriteSecrets(_ context.Context, path string, secrets map[string]string) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	if m.written == nil {
		m.written = make(map[string]map[string]string)
	}
	m.written[path] = secrets
	return nil
}

func newStore() *mockStore {
	return &mockStore{
		data: map[string]map[string]string{
			"secret/staging/app": {"DB_URL": "postgres://staging", "API_KEY": "abc"},
		},
	}
}

func TestPromote_DryRun_DoesNotWrite(t *testing.T) {
	store := newStore()
	p := promote.New(store, true, 5*time.Second)
	res, err := p.Promote(context.Background(), "secret/staging/app", "secret/prod/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
	if store.written != nil {
		t.Error("expected no writes in dry-run mode")
	}
}

func TestPromote_WritesSecrets(t *testing.T) {
	store := newStore()
	p := promote.New(store, false, 5*time.Second)
	res, err := p.Promote(context.Background(), "secret/staging/app", "secret/prod/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Keys))
	}
	if store.written["secret/prod/app"]["DB_URL"] != "postgres://staging" {
		t.Error("expected DB_URL to be promoted")
	}
}

func TestPromote_FetchError(t *testing.T) {
	store := &mockStore{fetchErr: errors.New("vault unavailable")}
	p := promote.New(store, false, 5*time.Second)
	_, err := p.Promote(context.Background(), "secret/staging/app", "secret/prod/app")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPromote_FullLifecycle(t *testing.T) {
	store := newStore()
	p := promote.New(store, false, 5*time.Second)
	res, err := p.Promote(context.Background(), "secret/staging/app", "secret/prod/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SourcePath != "secret/staging/app" {
		t.Errorf("unexpected source: %s", res.SourcePath)
	}
	if res.DestPath != "secret/prod/app" {
		t.Errorf("unexpected dest: %s", res.DestPath)
	}
}
