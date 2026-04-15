package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/diff"
)

func TestPrintDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	p.PrintDiff("staging.env", []diff.Change{})
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestPrintDiff_Added(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	changes := []diff.Change{
		{Key: "NEW_KEY", Type: diff.Added},
	}
	p.PrintDiff("prod.env", changes)
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added marker, got: %s", out)
	}
}

func TestPrintDiff_Removed(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	changes := []diff.Change{
		{Key: "OLD_KEY", Type: diff.Removed},
	}
	p.PrintDiff("prod.env", changes)
	out := buf.String()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removed marker, got: %s", out)
	}
}

func TestPrintDiff_Modified(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	changes := []diff.Change{
		{Key: "CHANGED_KEY", Type: diff.Modified},
	}
	p.PrintDiff("dev.env", changes)
	out := buf.String()
	if !strings.Contains(out, "~ CHANGED_KEY") {
		t.Errorf("expected modified marker, got: %s", out)
	}
}

// TestPrintDiff_IncludesFilename verifies that the filename is present in the
// diff output so users can identify which file is being reported.
func TestPrintDiff_IncludesFilename(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	changes := []diff.Change{
		{Key: "SOME_KEY", Type: diff.Added},
	}
	p.PrintDiff("myapp.env", changes)
	out := buf.String()
	if !strings.Contains(out, "myapp.env") {
		t.Errorf("expected filename 'myapp.env' in output, got: %s", out)
	}
}

func TestPrintSummary(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	changes := []diff.Change{
		{Key: "A", Type: diff.Added},
		{Key: "B", Type: diff.Removed},
		{Key: "C", Type: diff.Modified},
		{Key: "D", Type: diff.Unchanged},
	}
	p.PrintSummary("test.env", changes)
	out := buf.String()
	for _, want := range []string{"added: 1", "removed: 1", "modified: 1", "unchanged: 1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary output: %s", want, out)
		}
	}
}

func TestPrinter_ColorDisabled(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, true)
	if got := p.color(ColorGreen); got != "" {
		t.Errorf("expected empty string with noColor=true, got %q", got)
	}
}

func TestPrinter_ColorEnabled(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, false)
	if got := p.color(ColorGreen); got != string(ColorGreen) {
		t.Errorf("expected ANSI code, got %q", got)
	}
}
