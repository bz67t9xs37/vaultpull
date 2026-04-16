package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of requests per duration window.
type Limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	buckets  map[string][]time.Time
}

// Config holds rate limiter configuration.
type Config struct {
	Max    int
	Window time.Duration
}

// New creates a Limiter with the given config.
func New(cfg Config) *Limiter {
	if cfg.Max <= 0 {
		cfg.Max = 10
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	return &Limiter{
		max:     cfg.Max,
		window:  cfg.Window,
		buckets: make(map[string][]time.Time),
	}
}

// Allow reports whether the given key is within the rate limit.
// It records the attempt if allowed.
func (l *Limiter) Allow(key string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	var recent []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= l.max {
		return false, fmt.Errorf("rate limit exceeded for %q: %d requests in %s", key, l.max, l.window)
	}

	l.buckets[key] = append(recent, now)
	return true, nil
}

// Reset clears the rate limit state for a given key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Count returns the number of recent requests for a key within the window.
func (l *Limiter) Count(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
