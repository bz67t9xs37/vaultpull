package ttl

import (
	"fmt"
	"time"
)

// Entry holds a secret value along with its expiry metadata.
type Entry struct {
	Key       string
	Value     string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the entry has exceeded its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// ExpiresAt returns the absolute expiry time of the entry.
func (e *Entry) ExpiresAt() time.Time {
	return e.FetchedAt.Add(e.TTL)
}

// Tracker manages TTL entries for a set of secret keys.
type Tracker struct {
	entries map[string]*Entry
	defaultTTL time.Duration
}

// New creates a Tracker with the given default TTL.
func New(defaultTTL time.Duration) *Tracker {
	return &Tracker{
		entries:    make(map[string]*Entry),
		defaultTTL: defaultTTL,
	}
}

// Track registers a key with its value and the current time.
func (t *Tracker) Track(key, value string) {
	t.entries[key] = &Entry{
		Key:       key,
		Value:     value,
		FetchedAt: time.Now(),
		TTL:       t.defaultTTL,
	}
}

// Get returns the entry for a key, or nil if not tracked.
func (t *Tracker) Get(key string) *Entry {
	return t.entries[key]
}

// ExpiredKeys returns all keys whose TTL has elapsed.
func (t *Tracker) ExpiredKeys() []string {
	var expired []string
	for k, e := range t.entries {
		if e.IsExpired() {
			expired = append(expired, k)
		}
	}
	return expired
}

// Evict removes a key from the tracker.
func (t *Tracker) Evict(key string) {
	delete(t.entries, key)
}

// Summary returns a human-readable status line for a key.
func (t *Tracker) Summary(key string) string {
	e, ok := t.entries[key]
	if !ok {
		return fmt.Sprintf("%s: not tracked", key)
	}
	if e.TTL <= 0 {
		return fmt.Sprintf("%s: no TTL set", key)
	}
	if e.IsExpired() {
		return fmt.Sprintf("%s: expired (was due %s ago)", key, time.Since(e.ExpiresAt()).Truncate(time.Second))
	}
	return fmt.Sprintf("%s: valid for another %s", key, time.Until(e.ExpiresAt()).Truncate(time.Second))
}
