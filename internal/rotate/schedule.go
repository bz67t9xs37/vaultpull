package rotate

import (
	"fmt"
	"time"
)

// Schedule defines when rotation should occur.
type Schedule struct {
	// Interval is the duration between rotations.
	Interval time.Duration
	// LastRotated is the timestamp of the most recent rotation.
	LastRotated time.Time
}

// IsDue reports whether a rotation is due based on the schedule.
func (s *Schedule) IsDue() bool {
	if s.LastRotated.IsZero() {
		return true
	}
	return time.Since(s.LastRotated) >= s.Interval
}

// NextRotation returns the time at which the next rotation is due.
func (s *Schedule) NextRotation() time.Time {
	if s.LastRotated.IsZero() {
		return time.Now().UTC()
	}
	return s.LastRotated.Add(s.Interval)
}

// ParseInterval parses a human-readable duration string such as "24h" or "7d".
// It supports the standard Go duration suffixes plus "d" for days.
func ParseInterval(raw string) (time.Duration, error) {
	if len(raw) > 1 && raw[len(raw)-1] == 'd' {
		days := raw[:len(raw)-1]
		var n int
		if _, err := fmt.Sscanf(days, "%d", &n); err != nil {
			return 0, fmt.Errorf("invalid day value %q", raw)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid interval %q: %w", raw, err)
	}
	return d, nil
}
