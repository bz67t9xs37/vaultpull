package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Checker computes and compares checksums for secret maps.
type Checker struct {
	algorithm string
}

// Result holds the checksum result for a secret path.
type Result struct {
	Path     string
	Checksum string
	KeyCount int
}

// New returns a new Checker using SHA-256.
func New() *Checker {
	return &Checker{algorithm: "sha256"}
}

// Compute calculates a deterministic SHA-256 checksum over the given key-value map.
// Keys are sorted before hashing to ensure stability.
func (c *Checker) Compute(path string, secrets map[string]string) Result {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, secrets[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return Result{
		Path:     path,
		Checksum: hex.EncodeToString(sum[:]),
		KeyCount: len(secrets),
	}
}

// Equal returns true if two Results have the same checksum.
func (c *Checker) Equal(a, b Result) bool {
	return a.Checksum == b.Checksum
}

// Changed returns true if the checksum for the given secrets differs from the stored result.
func (c *Checker) Changed(previous Result, path string, current map[string]string) bool {
	next := c.Compute(path, current)
	return !c.Equal(previous, next)
}
