package snapshot

import (
	"sort"
)

// Drift describes the changes between two snapshots.
type Drift struct {
	// NewMissing are keys missing in the newer snapshot that were not missing before.
	NewMissing []string
	// ResolvedMissing are keys that were missing before but are no longer missing.
	ResolvedMissing []string
	// NewMismatched are keys that became mismatched between snapshots.
	NewMismatched []string
	// ResolvedMismatched are keys that were mismatched but are now clean.
	ResolvedMismatched []string
}

// IsClean returns true when there is no drift between the two snapshots.
func (d Drift) IsClean() bool {
	return len(d.NewMissing) == 0 &&
		len(d.ResolvedMissing) == 0 &&
		len(d.NewMismatched) == 0 &&
		len(d.ResolvedMismatched) == 0
}

// CompareDrift computes the drift between an older and a newer snapshot.
func CompareDrift(older, newer Snapshot) Drift {
	oldMissing := toSet(append(older.Result.MissingInRight, older.Result.MissingInLeft...))
	newMissing := toSet(append(newer.Result.MissingInRight, newer.Result.MissingInLeft...))

	oldMismatched := mismatchedKeys(older)
	newMismatched := mismatchedKeys(newer)

	drift := Drift{
		NewMissing:         setDiff(newMissing, oldMissing),
		ResolvedMissing:    setDiff(oldMissing, newMissing),
		NewMismatched:      setDiff(newMismatched, oldMismatched),
		ResolvedMismatched: setDiff(oldMismatched, newMismatched),
	}

	return drift
}

func mismatchedKeys(snap Snapshot) map[string]struct{} {
	m := make(map[string]struct{}, len(snap.Result.Mismatched))
	for _, mm := range snap.Result.Mismatched {
		m[mm.Key] = struct{}{}
	}
	return m
}

func toSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

func setDiff(a, b map[string]struct{}) []string {
	var result []string
	for k := range a {
		if _, found := b[k]; !found {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}
