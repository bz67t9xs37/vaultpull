package compare

import (
	"fmt"
	"sort"
)

// Result holds the outcome of comparing two secret maps.
type Result struct {
	Path    string
	Added   map[string]string
	Removed map[string]string
	Changed map[string]ChangedValue
	Equal   bool
}

// ChangedValue holds the before and after values for a modified key.
type ChangedValue struct {
	Before string
	After  string
}

// Comparer compares two maps of secrets.
type Comparer struct{}

// New returns a new Comparer.
func New() *Comparer {
	return &Comparer{}
}

// Compare returns a Result describing differences between src and dst.
func (c *Comparer) Compare(path string, src, dst map[string]string) Result {
	result := Result{
		Path:    path,
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]ChangedValue),
	}

	for k, v := range dst {
		if old, ok := src[k]; !ok {
			result.Added[k] = v
		} else if old != v {
			result.Changed[k] = ChangedValue{Before: old, After: v}
		}
	}

	for k, v := range src {
		if _, ok := dst[k]; !ok {
			result.Removed[k] = v
		}
	}

	result.Equal = len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0
	return result
}

// Summary returns a human-readable summary of the result.
func (r Result) Summary() string {
	if r.Equal {
		return fmt.Sprintf("%s: no differences", r.Path)
	}
	return fmt.Sprintf("%s: +%d added, -%d removed, ~%d changed",
		r.Path, len(r.Added), len(r.Removed), len(r.Changed))
}

// Keys returns a sorted list of all affected keys.
func (r Result) Keys() []string {
	seen := make(map[string]struct{})
	for k := range r.Added {
		seen[k] = struct{}{}
	}
	for k := range r.Removed {
		seen[k] = struct{}{}
	}
	for k := range r.Changed {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
