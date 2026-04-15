package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/diff"
	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Path    string
	Summary diff.Summary
	Diff    []diff.Change
}

// Syncer orchestrates pulling secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
}

// New creates a new Syncer with the given Vault client.
func New(client *vault.Client) *Syncer {
	return &Syncer{client: client}
}

// Sync fetches secrets at secretPath from Vault and merges them into the .env
// file at envPath. It returns a Result describing what changed.
func (s *Syncer) Sync(secretPath, envPath string) (*Result, error) {
	remote, err := s.client.GetSecrets(secretPath)
	if err != nil {
		return nil, fmt.Errorf("fetching secrets from vault path %q: %w", secretPath, err)
	}

	local, err := envfile.Parse(envPath)
	if err != nil {
		return nil, fmt.Errorf("parsing env file %q: %w", envPath, err)
	}

	changes := diff.Compute(local, remote)
	summary := diff.Summary(changes)

	if summary.Added == 0 && summary.Modified == 0 && summary.Removed == 0 {
		return &Result{
			Path:    envPath,
			Summary: summary,
			Diff:    changes,
		}, nil
	}

	if err := envfile.Write(envPath, remote); err != nil {
		return nil, fmt.Errorf("writing env file %q: %w", envPath, err)
	}

	return &Result{
		Path:    envPath,
		Summary: summary,
		Diff:    changes,
	}, nil
}
