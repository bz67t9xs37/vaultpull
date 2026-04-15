package sync

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultpull/internal/backup"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/output"
)

// VaultClient fetches secrets from Vault.
type VaultClient interface {
	GetSecrets(ctx context.Context, mount, path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing .env files.
type Syncer struct {
	client  VaultClient
	printer *output.Printer
	backups *backup.Manager
}

// New creates a Syncer wired with the given client and printer.
// If cfg.BackupDir is set, backup support is enabled.
func New(client VaultClient, printer *output.Printer, cfg *config.Config) *Syncer {
	var bm *backup.Manager
	if cfg.BackupDir != "" {
		bm = backup.New(cfg.BackupDir)
	}
	return &Syncer{client: client, printer: printer, backups: bm}
}

// Run iterates over all targets, pulls secrets, diffs, and writes.
func (s *Syncer) Run(ctx context.Context, cfg *config.Config) error {
	for _, target := range cfg.Targets {
		if err := s.syncTarget(ctx, cfg.Mount, target); err != nil {
			return fmt.Errorf("sync target %q: %w", target.Path, err)
		}
	}
	return nil
}

func (s *Syncer) syncTarget(ctx context.Context, mount string, target config.Target) error {
	secrets, err := s.client.GetSecrets(ctx, mount, target.Path)
	if err != nil {
		return err
	}

	existing, err := envfile.Parse(target.Output)
	if err != nil {
		return err
	}

	changes := diff.Compute(existing, secrets)
	s.printer.PrintDiff(target.Output, changes)
	s.printer.PrintSummary(diff.Summary(changes))

	if diff.Summary(changes).Total() == 0 {
		return nil
	}

	if target.Backup && s.backups != nil {
		if _, err := s.backups.Create(target.Output); err != nil {
			return fmt.Errorf("create backup: %w", err)
		}
	}

	return envfile.Write(target.Output, secrets)
}
