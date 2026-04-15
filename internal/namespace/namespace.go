package namespace

import (
	"fmt"
	"strings"
)

// Resolver maps logical target names to Vault secret paths using
// a configurable namespace prefix.
type Resolver struct {
	prefix string
	mount  string
}

// New creates a new Resolver with the given namespace prefix and KV mount.
func New(prefix, mount string) *Resolver {
	return &Resolver{
		prefix: strings.Trim(prefix, "/"),
		mount:  strings.Trim(mount, "/"),
	}
}

// Resolve returns the full Vault path for a given secret name.
// Example: prefix="myapp/prod", mount="secret", name="db" => "secret/data/myapp/prod/db"
func (r *Resolver) Resolve(name string) string {
	name = strings.Trim(name, "/")
	if r.prefix == "" {
		return fmt.Sprintf("%s/data/%s", r.mount, name)
	}
	return fmt.Sprintf("%s/data/%s/%s", r.mount, r.prefix, name)
}

// ResolveAll returns full Vault paths for a list of secret names.
func (r *Resolver) ResolveAll(names []string) []string {
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = r.Resolve(name)
	}
	return paths
}

// StripMount removes the mount and "/data/" prefix from a full Vault path,
// returning the logical secret path relative to the namespace.
func (r *Resolver) StripMount(fullPath string) string {
	prefix := r.mount + "/data/"
	result := strings.TrimPrefix(fullPath, prefix)
	if r.prefix != "" {
		result = strings.TrimPrefix(result, r.prefix+"/")
	}
	return result
}
