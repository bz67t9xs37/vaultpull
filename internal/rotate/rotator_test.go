package rotate_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/backup"
	"github.com/user/vaultpull/internal/rotate"
	"github.com/user/vaultpull/internal/vault"
)

func newMockVault(t *testing.T, secrets map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := `{"data":{"data":{`
		parts := []string{}
		for k, v := range secrets {
			parts = append(parts, `"`+k+`":"`+v+`"`)
		}
		data += strings.Join(parts, ",") + `}}}`
		w.Write([]byte(data))
	}))
}

func TestRotate_Success(t *testing.T) {
	srv := newMockVault(t, map[string]string{"API_KEY": "newval"})
	defer srv.Close()

	vc, err := vault.NewClient(srv.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("API_KEY=oldval\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	bs := backup.New(tmpDir)
	rot := rotate.New(vc, bs)

	res, err := rot.Rotate(envPath, "secret", "myapp")
	if err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	if res.BackupPath == "" {
		t.Error("expected backup path to be set")
	}
	if len(res.KeysRotated) == 0 {
		t.Error("expected at least one rotated key")
	}
	if res.RotatedAt.IsZero() {
		t.Error("expected RotatedAt to be set")
	}
}

func TestRotate_MissingEnvFile(t *testing.T) {
	srv := newMockVault(t, map[string]string{"X": "y"})
	defer srv.Close()

	vc, _ := vault.NewClient(srv.URL, "tok", "secret")
	tmpDir := t.TempDir()
	bs := backup.New(tmpDir)
	rot := rotate.New(vc, bs)

	// env file does not exist — backup should fail gracefully
	_, err := rot.Rotate(filepath.Join(tmpDir, "nonexistent.env"), "secret", "myapp")
	if err == nil {
		t.Fatal("expected error for missing env file")
	}
}
