package rollback

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single rollback point for a target env file.
type Entry struct {
	Path      string
	BackupAt  time.Time
	BackupFile string
}

// Rollbacker restores env files from backup snapshots.
type Rollbacker struct {
	backupDir string
}

// New creates a Rollbacker that looks for backups in backupDir.
func New(backupDir string) *Rollbacker {
	return &Rollbacker{backupDir: backupDir}
}

// Latest returns the most recent backup entry for the given target path,
// or nil if no backup exists.
func (r *Rollbacker) Latest(targetPath string) (*Entry, error) {
	pattern := filepath.Join(r.backupDir, backupGlob(targetPath))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("rollback: glob backups: %w", err)
	}
	if len(matches) == 0 {
		return nil, nil
	}
	latest := matches[len(matches)-1]
	info, err := os.Stat(latest)
	if err != nil {
		return nil, fmt.Errorf("rollback: stat backup: %w", err)
	}
	return &Entry{
		Path:       targetPath,
		BackupAt:   info.ModTime(),
		BackupFile: latest,
	}, nil
}

// Restore copies the latest backup over the target path.
func (r *Rollbacker) Restore(targetPath string) (*Entry, error) {
	entry, err := r.Latest(targetPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("rollback: no backup found for %s", targetPath)
	}
	data, err := os.ReadFile(entry.BackupFile)
	if err != nil {
		return nil, fmt.Errorf("rollback: read backup: %w", err)
	}
	if err := os.WriteFile(targetPath, data, 0o600); err != nil {
		return nil, fmt.Errorf("rollback: write target: %w", err)
	}
	return entry, nil
}

// backupGlob returns a glob pattern matching backup files for a given path.
func backupGlob(targetPath string) string {
	base := filepath.Base(targetPath)
	return base + ".*.bak"
}
