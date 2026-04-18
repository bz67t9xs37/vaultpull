package promote

import (
	"context"
	"fmt"
	"time"
)

// SecretStore abstracts reading and writing secrets.
type SecretStore interface {
	GetSecrets(ctx context.Context, path string) (map[string]string, error)
	WriteSecrets(ctx context.Context, path string, secrets map[string]string) error
}

// Promoter copies secrets from a source path to a destination path.
type Promoter struct {
	store   SecretStore
	dryRun  bool
	timeout time.Duration
}

// New creates a new Promoter.
func New(store SecretStore, dryRun bool, timeout time.Duration) *Promoter {
	return &Promoter{store: store, dryRun: dryRun, timeout: timeout}
}

// Result holds the outcome of a promotion.
type Result struct {
	SourcePath string
	DestPath   string
	Keys       []string
	DryRun     bool
}

// Promote fetches secrets from src and writes them to dst.
func (p *Promoter) Promote(ctx context.Context, src, dst string) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	secrets, err := p.store.GetSecrets(ctx, src)
	if err != nil {
		return nil, fmt.Errorf("promote: fetch from %q: %w", src, err)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	if !p.dryRun {
		if err := p.store.WriteSecrets(ctx, dst, secrets); err != nil {
			return nil, fmt.Errorf("promote: write to %q: %w", dst, err)
		}
	}

	return &Result{
		SourcePath: src,
		DestPath:   dst,
		Keys:       keys,
		DryRun:     p.dryRun,
	}, nil
}
