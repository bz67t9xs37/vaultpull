package rotate

import (
	"fmt"
	"time"

	"github.com/user/vaultpull/internal/backup"
	"github.com/user/vaultpull/internal/envfile"
	"github.com/user/vaultpull/internal/vault"
)

// Result holds the outcome of a rotation operation.
type Result struct {
	Target    string
	BackupPath string
	RotatedAt time.Time
	KeysRotated []string
}

// Rotator handles secret rotation for a single target.
type Rotator struct {
	vaultClient *vault.Client
	backupSvc   *backup.Service
}

// New creates a new Rotator.
func New(vc *vault.Client, bs *backup.Service) *Rotator {
	return &Rotator{
		vaultClient: vc,
		backupSvc:   bs,
	}
}

// Rotate fetches fresh secrets from Vault, backs up the existing env file,
// and overwrites it with the latest values. It returns a Result describing
// what changed.
func (r *Rotator) Rotate(envPath, mountPath, secretPath string) (*Result, error) {
	// Back up existing file before overwriting.
	bp, err := r.backupSvc.Create(envPath)
	if err != nil {
		return nil, fmt.Errorf("rotate: backup failed: %w", err)
	}

	// Fetch latest secrets from Vault.
	secrets, err := r.vaultClient.GetSecrets(mountPath, secretPath)
	if err != nil {
		return nil, fmt.Errorf("rotate: vault fetch failed: %w", err)
	}

	// Write new secrets to env file.
	if err := envfile.Write(envPath, secrets); err != nil {
		return nil, fmt.Errorf("rotate: write failed: %w", err)
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	return &Result{
		Target:      envPath,
		BackupPath:  bp,
		RotatedAt:   time.Now().UTC(),
		KeysRotated: keys,
	}, nil
}
