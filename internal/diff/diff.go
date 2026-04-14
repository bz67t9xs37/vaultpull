package diff

import "fmt"

// ChangeType represents the type of change for a secret entry.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single diff entry between old and new secret values.
type Change struct {
	Key    string
	OldVal string
	NewVal string
	Type   ChangeType
}

// Compute returns the diff between existing (old) and incoming (new) secret maps.
func Compute(old, new map[string]string) []Change {
	var changes []Change

	for k, newVal := range new {
		if oldVal, exists := old[k]; !exists {
			changes = append(changes, Change{Key: k, OldVal: "", NewVal: newVal, Type: Added})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: newVal, Type: Modified})
		} else {
			changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: newVal, Type: Unchanged})
		}
	}

	for k, oldVal := range old {
		if _, exists := new[k]; !exists {
			changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: "", Type: Removed})
		}
	}

	return changes
}

// Summary returns a human-readable summary string of the diff.
func Summary(changes []Change) string {
	added, removed, modified := 0, 0, 0
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d modified", added, removed, modified)
}
