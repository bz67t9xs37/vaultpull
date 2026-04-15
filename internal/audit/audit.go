package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Target    string    `json:"target"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Modified  int       `json:"modified"`
	Unchanged int       `json:"unchanged"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a JSON-lines log file.
type Logger struct {
	path string
}

// New creates a Logger that appends to the file at path.
func New(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	return &Logger{path: path}, nil
}

// Record appends an Entry to the audit log file.
func (l *Logger) Record(e Entry) error {
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// ReadAll parses and returns all entries from the audit log file.
func (l *Logger) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
