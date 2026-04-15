package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles creating and restoring .env file backups.
type Manager struct {
	backupDir string
}

// New creates a new backup Manager storing backups in backupDir.
func New(backupDir string) *Manager {
	return &Manager{backupDir: backupDir}
}

// Create writes a timestamped backup of the file at src.
// Returns the path of the created backup file.
func (m *Manager) Create(src string) (string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("backup: read source %q: %w", src, err)
	}

	if err := os.MkdirAll(m.backupDir, 0o700); err != nil {
		return "", fmt.Errorf("backup: create dir %q: %w", m.backupDir, err)
	}

	base := filepath.Base(src)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dest := filepath.Join(m.backupDir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return "", fmt.Errorf("backup: write backup %q: %w", dest, err)
	}

	return dest, nil
}

// Restore copies the backup file at backupPath back to dest.
func (m *Manager) Restore(backupPath, dest string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("backup: read backup %q: %w", backupPath, err)
	}

	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("backup: restore to %q: %w", dest, err)
	}

	return nil
}

// List returns all backup files for a given original filename.
func (m *Manager) List(originalName string) ([]string, error) {
	pattern := filepath.Join(m.backupDir, originalName+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("backup: list %q: %w", pattern, err)
	}
	return matches, nil
}
