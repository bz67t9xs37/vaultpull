package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/template"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "template.env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempTemplate: %v", err)
	}
	return p
}

func TestRender_SubstitutesKnownKeys(t *testing.T) {
	r := template.New()
	src := "DB_HOST=${DB_HOST}\nDB_PASS=${DB_PASS}"
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASS": "s3cr3t"}
	got := r.Render(src, secrets)
	want := "DB_HOST=localhost\nDB_PASS=s3cr3t"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRender_LeavesUnknownPlaceholders(t *testing.T) {
	r := template.New()
	src := "KEY=${UNKNOWN}"
	got := r.Render(src, map[string]string{})
	if got != src {
		t.Errorf("Render() = %q, want original %q", got, src)
	}
}

func TestRenderFile_ReadsAndSubstitutes(t *testing.T) {
	r := template.New()
	p := writeTempTemplate(t, "TOKEN=${TOKEN}")
	got, err := r.RenderFile(p, map[string]string{"TOKEN": "abc123"})
	if err != nil {
		t.Fatalf("RenderFile() error: %v", err)
	}
	if got != "TOKEN=abc123" {
		t.Errorf("RenderFile() = %q", got)
	}
}

func TestRenderFile_MissingFile(t *testing.T) {
	r := template.New()
	_, err := r.RenderFile("/no/such/file.env", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestListPlaceholders_Unique(t *testing.T) {
	r := template.New()
	src := "${FOO}=${FOO} ${BAR}"
	keys := r.ListPlaceholders(src)
	if len(keys) != 2 {
		t.Errorf("ListPlaceholders() = %v, want 2 unique keys", keys)
	}
}

func TestMissingKeys_ReturnsMissing(t *testing.T) {
	r := template.New()
	src := "A=${A}\nB=${B}\nC=${C}"
	secrets := map[string]string{"A": "1"}
	missing := r.MissingKeys(src, secrets)
	if len(missing) != 2 {
		t.Errorf("MissingKeys() = %v, want 2 missing", missing)
	}
}

func TestValidateTemplate_ReturnsErrorOnMissing(t *testing.T) {
	r := template.New()
	p := writeTempTemplate(t, "SECRET=${MISSING_KEY}")
	if err := r.ValidateTemplate(p, map[string]string{}); err == nil {
		t.Fatal("expected validation error for unresolved placeholder")
	}
}

func TestValidateTemplate_PassesWhenAllResolved(t *testing.T) {
	r := template.New()
	p := writeTempTemplate(t, "SECRET=${PRESENT}")
	if err := r.ValidateTemplate(p, map[string]string{"PRESENT": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
