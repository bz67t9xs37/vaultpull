package promote_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/promote"
)

type inMemoryStore struct {
	data map[string]map[string]string
}

func newStore() *inMemoryStore {
	return &inMemoryStore{data: make(map[string]map[string]string)}
}

func (s *inMemoryStore) GetSecrets(_ context.Context, path string) (map[string]string, error) {
	if v, ok := s.data[path]; ok {
		return v, nil
	}
	return map[string]string{}, nil
}

func (s *inMemoryStore) WriteSecrets(_ context.Context, path string, secrets map[string]string) error {
	s.data[path] = secrets
	return nil
}

func TestPromote_FullLifecycle(t *testing.T) {
	store := newStore()
	store.data["secret/staging"] = map[string]string{
		"API_KEY": "staging-key",
		"DB_URL":  "postgres://staging",
	}

	p := promote.New(store, store, 5*time.Second, false)
	result, err := p.Promote("secret/staging", "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Keys) != 2 {
		t.Errorf("expected 2 promoted keys, got %d", len(result.Keys))
	}
	if store.data["secret/prod"]["API_KEY"] != "staging-key" {
		t.Error("expected prod to have promoted API_KEY")
	}
}

func TestPromote_DryRun_LeavesDestinationEmpty(t *testing.T) {
	store := newStore()
	store.data["secret/staging"] = map[string]string{"TOKEN": "abc"}

	p := promote.New(store, store, 5*time.Second, true)
	_, err := p.Promote("secret/staging", "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := store.data["secret/prod"]; exists {
		t.Error("dry run should not write to destination")
	}
}
