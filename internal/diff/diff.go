package diff

// ChangeType describes the kind of change for a key.
type ChangeType int

const (
	Added     ChangeType = iota
	Removed
	Modified
	Unchanged
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Type     ChangeType
}

// SummaryResult holds counts of each change type.
type SummaryResult struct {
	Added     int
	Removed   int
	Modified  int
	Unchanged int
}

// Compute compares two maps and returns a slice of Changes.
func Compute(existing, incoming map[string]string) []Change {
	var changes []Change

	for k, newVal := range incoming {
		if oldVal, ok := existing[k]; !ok {
			changes = append(changes, Change{Key: k, NewValue: newVal, Type: Added})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Type: Modified})
		} else {
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Type: Unchanged})
		}
	}

	for k, oldVal := range existing {
		if _, ok := incoming[k]; !ok {
			changes = append(changes, Change{Key: k, OldValue: oldVal, Type: Removed})
		}
	}

	return changes
}

// Summary aggregates a slice of Changes into a SummaryResult.
func Summary(changes []Change) SummaryResult {
	var s SummaryResult
	for _, c := range changes {
		switch c.Type {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		case Unchanged:
			s.Unchanged++
		}
	}
	return s
}
