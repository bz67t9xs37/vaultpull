package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpull/internal/diff"
)

// ColorCode represents an ANSI color escape code.
type ColorCode string

const (
	ColorReset  ColorCode = "\033[0m"
	ColorRed    ColorCode = "\033[31m"
	ColorGreen  ColorCode = "\033[32m"
	ColorYellow ColorCode = "\033[33m"
	ColorCyan   ColorCode = "\033[36m"
)

// Printer writes sync results to an output stream.
type Printer struct {
	w      io.Writer
	noColor bool
}

// New creates a Printer writing to w. Pass noColor=true to disable ANSI codes.
func New(w io.Writer, noColor bool) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w, noColor: noColor}
}

// PrintDiff writes a human-readable diff to the printer's writer.
func (p *Printer) PrintDiff(target string, changes []diff.Change) {
	if len(changes) == 0 {
		fmt.Fprintf(p.w, "%s: no changes\n", target)
		return
	}
	fmt.Fprintf(p.w, "%s:\n", target)
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			fmt.Fprintf(p.w, "  %s+ %s%s\n", p.color(ColorGreen), c.Key, p.color(ColorReset))
		case diff.Removed:
			fmt.Fprintf(p.w, "  %s- %s%s\n", p.color(ColorRed), c.Key, p.color(ColorReset))
		case diff.Modified:
			fmt.Fprintf(p.w, "  %s~ %s%s\n", p.color(ColorYellow), c.Key, p.color(ColorReset))
		case diff.Unchanged:
			fmt.Fprintf(p.w, "    %s\n", c.Key)
		}
	}
}

// PrintSummary writes a one-line summary of the diff.
func (p *Printer) PrintSummary(target string, changes []diff.Change) {
	s := diff.Summary(changes)
	fmt.Fprintf(p.w, "%s%s%s — added: %d, removed: %d, modified: %d, unchanged: %d\n",
		p.color(ColorCyan), target, p.color(ColorReset),
		s.Added, s.Removed, s.Modified, s.Unchanged)
}

func (p *Printer) color(c ColorCode) string {
	if p.noColor {
		return ""
	}
	return string(c)
}
