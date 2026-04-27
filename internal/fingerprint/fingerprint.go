package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Entry holds a computed fingerprint for a secret path at a point in time.
type Entry struct {
	Path      string
	Hash      string
	KeyCount  int
	ComputedAt time.Time
}

// Fingerprinter computes and compares fingerprints for secret maps.
type Fingerprinter struct {
	clock func() time.Time
}

// New returns a Fingerprinter using the real clock.
func New() *Fingerprinter {
	return &Fingerprinter{clock: time.Now}
}

func newWithClock(clock func() time.Time) *Fingerprinter {
	return &Fingerprinter{clock: clock}
}

// Compute returns an Entry for the given path and secret map.
func (f *Fingerprinter) Compute(path string, secrets map[string]string) Entry {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}

	return Entry{
		Path:       path,
		Hash:       hex.EncodeToString(h.Sum(nil)),
		KeyCount:   len(secrets),
		ComputedAt: f.clock(),
	}
}

// Equal returns true if two entries share the same hash.
func Equal(a, b Entry) bool {
	return strings.EqualFold(a.Hash, b.Hash)
}

// Changed returns true if the current secrets produce a different hash than the stored entry.
func (f *Fingerprinter) Changed(stored Entry, path string, current map[string]string) bool {
	next := f.Compute(path, current)
	return !Equal(stored, next)
}
