package drift

import (
	"fmt"
	"time"
)

// Result holds the outcome of a drift detection check.
type Result struct {
	Path      string
	Key       string
	VaultValue string
	LocalValue string
	Drifted   bool
	CheckedAt time.Time
}

// Detector compares vault secrets against local env values.
type Detector struct {
	secrets map[string]string
}

// New creates a Detector seeded with the current vault secrets.
func New(vaultSecrets map[string]string) *Detector {
	return &Detector{secrets: vaultSecrets}
}

// Check compares localEnv against the vault secrets and returns drift results.
func (d *Detector) Check(path string, localEnv map[string]string) []Result {
	now := time.Now().UTC()
	var results []Result
	for key, vaultVal := range d.secrets {
		localVal, exists := localEnv[key]
		drifted := !exists || localVal != vaultVal
		results = append(results, Result{
			Path:       path,
			Key:        key,
			VaultValue: vaultVal,
			LocalValue: localVal,
			Drifted:    drifted,
			CheckedAt:  now,
		})
	}
	return results
}

// Summary returns a human-readable summary of drift results.
func Summary(results []Result) string {
	total, drifted := 0, 0
	for _, r := range results {
		total++
		if r.Drifted {
			drifted++
		}
	}
	if drifted == 0 {
		return fmt.Sprintf("no drift detected (%d keys checked)", total)
	}
	return fmt.Sprintf("%d/%d keys drifted", drifted, total)
}
