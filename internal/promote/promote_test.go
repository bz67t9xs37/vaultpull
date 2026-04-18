package promote

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockFetcher struct {
	secrets map[string]string
	err     error
}

func (m *mockFetcher) GetSecrets(_ context.Context, _ string) (map[string]string, error) {
	return m.secrets, m.err
}

type mockWriter struct {
	written map[string]string
	err     error
}

func (m *mockWriter) WriteSecrets(_ context.Context, _ string, secrets map[string]string) error {
	m.written = secrets
	return m.err
}

func TestPromote_DryRun_DoesNotWrite(t *testing.T) {
	fetcher := &mockFetcher{secrets: map[string]string{"KEY": "val"}}
	writer := &mockWriter{}
	p := New(fetcher, writer, 5*time.Second, true)

	result, err := p.Promote("secret/staging", "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected dry run")
	}
	if writer.written != nil {
		t.Error("expected no write in dry run")
	}
	if len(result.Keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(result.Keys))
	}
}

func TestPromote_WritesSecrets(t *testing.T) {
	fetcher := &mockFetcher{secrets: map[string]string{"DB_PASS": "secret"}}
	writer := &mockWriter{}
	p := New(fetcher, writer, 5*time.Second, false)

	result, err := p.Promote("secret/staging", "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.DryRun {
		t.Error("expected live run")
	}
	if writer.written["DB_PASS"] != "secret" {
		t.Errorf("expected written secret, got %v", writer.written)
	}
}

func TestPromote_FetchError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("vault unavailable")}
	writer := &mockWriter{}
	p := New(fetcher, writer, 5*time.Second, false)

	_, err := p.Promote("secret/staging", "secret/prod")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPromote_WriteError(t *testing.T) {
	fetcher := &mockFetcher{secrets: map[string]string{"X": "y"}}
	writer := &mockWriter{err: errors.New("write failed")}
	p := New(fetcher, writer, 5*time.Second, false)

	_, err := p.Promote("secret/staging", "secret/prod")
	if err == nil {
		t.Fatal("expected error")
	}
}
