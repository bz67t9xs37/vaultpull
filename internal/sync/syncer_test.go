package sync_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

func newMockVault(t *testing.T, secrets map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := `{"data":{"data":{`
		i := 0
		for k, v := range secrets {
			if i > 0 {
				data += ","
			}
			data += `"` + k + `":"` + v + `"`
			i++
		}
		data += `}}}`
		w.Write([]byte(data))
	}))
}

func TestSync_NewSecrets(t *testing.T) {
	server := newMockVault(t, map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	s := sync.New(client)
	result, err := s.Sync("myapp/config", envPath)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Summary.Added != 2 {
		t.Errorf("expected 2 added, got %d", result.Summary.Added)
	}
	if result.Summary.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", result.Summary.Modified)
	}

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error("expected env file to be created")
	}
}

func TestSync_NoChanges(t *testing.T) {
	server := newMockVault(t, map[string]string{"FOO": "bar"})
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("FOO=bar\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	s := sync.New(client)
	result, err := s.Sync("myapp/config", envPath)
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if result.Summary.Added != 0 || result.Summary.Modified != 0 || result.Summary.Removed != 0 {
		t.Errorf("expected no changes, got %+v", result.Summary)
	}
}
