package promote

import (
	"context"
	"fmt"
	"time"
)

// SecretFetcher retrieves secrets from a given path.
type SecretFetcher interface {
	GetSecrets(ctx context.Context, path string) (map[string]string, error)
}

// SecretWriter writes secrets to a given path.
type SecretWriter interface {
	WriteSecrets(ctx context.Context, path string, secrets map[string]string) error
}

// Result holds the outcome of a promotion.
type Result struct {
	Source      string
	Destination string
	Keys        []string
	DryRun      bool
	PromotedAt  time.Time
}

// Promoter copies secrets from one path to another.
type Promoter struct {
	fetcher SecretFetcher
	writer  SecretWriter
	timeout time.Duration
	dryRun  bool
}

// New creates a new Promoter.
func New(fetcher SecretFetcher, writer SecretWriter, timeout time.Duration, dryRun bool) *Promoter {
	return &Promoter{
		fetcher: fetcher,
		writer:  writer,
		timeout: timeout,
		dryRun:  dryRun,
	}
}

// Promote copies secrets from source to destination path.
func (p *Promoter) Promote(source, destination string) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	secrets, err := p.fetcher.GetSecrets(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("promote: fetch from %q: %w", source, err)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	if !p.dryRun {
		if err := p.writer.WriteSecrets(ctx, destination, secrets); err != nil {
			return nil, fmt.Errorf("promote: write to %q: %w", destination, err)
		}
	}

	return &Result{
		Source:      source,
		Destination: destination,
		Keys:        keys,
		DryRun:      p.dryRun,
		PromotedAt:  time.Now(),
	}, nil
}
