package resolve

import (
	"fmt"
	"strings"
)

// Resolver resolves secret paths from a mount and a list of keys,
// supporting both explicit paths and glob-style prefix expansion.
type Resolver struct {
	mount string
}

// ResolvedPath holds the full Vault path and the local key name.
type ResolvedPath struct {
	VaultPath string
	LocalKey  string
}

// New creates a new Resolver for the given mount point.
func New(mount string) *Resolver {
	return &Resolver{
		mount: strings.Trim(mount, "/"),
	}
}

// Resolve converts a list of secret keys into full Vault paths.
// Each key may be a simple name (e.g. "db/password") or an absolute
// path already containing the mount (e.g. "secret/db/password").
func (r *Resolver) Resolve(keys []string) ([]ResolvedPath, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("resolve: no keys provided")
	}

	paths := make([]ResolvedPath, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		var vaultPath string
		if strings.HasPrefix(key, r.mount+"/") {
			// Already fully qualified.
			vaultPath = key
		} else {
			vaultPath = r.mount + "/" + strings.TrimPrefix(key, "/")
		}

		paths = append(paths, ResolvedPath{
			VaultPath: vaultPath,
			LocalKey:  localKeyName(key),
		})
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("resolve: all keys were empty after trimming")
	}

	return paths, nil
}

// Mount returns the configured mount point.
func (r *Resolver) Mount() string {
	return r.mount
}

// localKeyName derives a safe local environment-variable key from a path.
// e.g. "secret/db/password" -> "DB_PASSWORD"
func localKeyName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Drop the mount segment if it looks like a generic mount name.
	if len(parts) > 1 {
		parts = parts[1:]
	}
	joined := strings.Join(parts, "_")
	return strings.ToUpper(strings.ReplaceAll(joined, "-", "_"))
}
