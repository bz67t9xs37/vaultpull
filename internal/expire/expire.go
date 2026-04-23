package expire

import (
	"fmt"
	"time"
)

// Entry holds expiry metadata for a single secret path.
type Entry struct {
	Path      string
	ExpiresAt time.Time
	Warned    bool
}

// Checker evaluates whether secrets are approaching or past expiry.
type Checker struct {
	warnBefore time.Duration
	clock      func() time.Time
}

// Result represents the expiry status of a single path.
type Result struct {
	Path    string
	Expired bool
	Warning bool
	// TimeLeft is negative when already expired.
	TimeLeft time.Duration
}

// New returns a Checker. warnBefore is the window before expiry to emit warnings.
func New(warnBefore time.Duration) *Checker {
	return &Checker{
		warnBefore: warnBefore,
		clock:      time.Now,
	}
}

// Check evaluates a slice of entries and returns results for each.
func (c *Checker) Check(entries []Entry) []Result {
	now := c.clock()
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		if e.ExpiresAt.IsZero() {
			continue
		}
		timeLeft := e.ExpiresAt.Sub(now)
		results = append(results, Result{
			Path:     e.Path,
			Expired:  timeLeft <= 0,
			Warning:  timeLeft > 0 && timeLeft <= c.warnBefore,
			TimeLeft: timeLeft,
		})
	}
	return results
}

// Summary returns a human-readable summary of expiry results.
func Summary(results []Result) string {
	expired, warned := 0, 0
	for _, r := range results {
		if r.Expired {
			expired++
		} else if r.Warning {
			warned++
		}
	}
	if expired == 0 && warned == 0 {
		return "all secrets are within their expiry window"
	}
	return fmt.Sprintf("%d expired, %d expiring soon", expired, warned)
}
