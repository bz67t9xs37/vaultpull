package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpull/internal/vault"
)

func newMockVaultServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			http.Error(w, "missing token", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestGetSecrets_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
		},
	}

	srv := newMockVaultServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	secrets, err := client.GetSecrets(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("GetSecrets() error = %v", err)
	}

	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
	})
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestNewClient_DefaultMount(t *testing.T) {
	client, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
		Token:   "tok",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.Mount != "secret" {
		t.Errorf("expected default mount=secret, got %q", client.Mount)
	}
}
