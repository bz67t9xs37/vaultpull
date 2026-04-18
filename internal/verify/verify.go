package verify

import (
	"errors"
	"fmt"
)

// Result holds the outcome of verifying a single secret key.
type Result struct {
	Key     string
	Present bool
	Empty   bool
	Err     error
}

// Report aggregates verification results for a path.
type Report struct {
	Path    string
	Results []Result
}

// Summary returns a human-readable summary of the report.
func (r *Report) Summary() string {
	missing, empty := 0, 0
	for _, res := range r.Results {
		if !res.Present {
			missing++
		} else if res.Empty {
			empty++
		}
	}
	if missing == 0 && empty == 0 {
		return fmt.Sprintf("%s: all %d keys verified", r.Path, len(r.Results))
	}
	return fmt.Sprintf("%s: %d missing, %d empty", r.Path, missing, empty)
}

// HasErrors returns true if any result is missing or errored.
func (r *Report) HasErrors() bool {
	for _, res := range r.Results {
		if !res.Present || res.Err != nil {
			return true
		}
	}
	return false
}

// Verifier checks that expected keys exist and are non-empty in a secrets map.
type Verifier struct {
	requireNonEmpty bool
}

// New returns a Verifier. If requireNonEmpty is true, empty values are flagged.
func New(requireNonEmpty bool) *Verifier {
	return &Verifier{requireNonEmpty: requireNonEmpty}
}

// Check verifies that all expectedKeys are present in secrets.
func (v *Verifier) Check(path string, secrets map[string]string, expectedKeys []string) (*Report, error) {
	if len(expectedKeys) == 0 {
		return nil, errors.New("verify: no expected keys provided")
	}
	report := &Report{Path: path}
	for _, key := range expectedKeys {
		val, ok := secrets[key]
		res := Result{Key: key, Present: ok}
		if ok && v.requireNonEmpty && val == "" {
			res.Empty = true
		}
		report.Results = append(report.Results, res)
	}
	return report, nil
}
