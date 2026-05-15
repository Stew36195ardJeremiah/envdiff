package snapshot_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func makeSnap(missingRight, missingLeft []string, mismatched []diff.Mismatch) snapshot.Snapshot {
	return snapshot.Snapshot{
		Result: diff.Result{
			MissingInRight: missingRight,
			MissingInLeft:  missingLeft,
			Mismatched:     mismatched,
		},
	}
}

func TestCompareDrift_Clean(t *testing.T) {
	old := makeSnap(nil, nil, nil)
	new := makeSnap(nil, nil, nil)
	drift := snapshot.CompareDrift(old, new)
	if !drift.IsClean() {
		t.Errorf("expected clean drift, got %+v", drift)
	}
}

func TestCompareDrift_NewMissing(t *testing.T) {
	old := makeSnap(nil, nil, nil)
	new := makeSnap([]string{"DB_HOST"}, nil, nil)
	drift := snapshot.CompareDrift(old, new)
	if len(drift.NewMissing) != 1 || drift.NewMissing[0] != "DB_HOST" {
		t.Errorf("expected NewMissing=[DB_HOST], got %v", drift.NewMissing)
	}
}

func TestCompareDrift_ResolvedMissing(t *testing.T) {
	old := makeSnap([]string{"DB_HOST"}, nil, nil)
	new := makeSnap(nil, nil, nil)
	drift := snapshot.CompareDrift(old, new)
	if len(drift.ResolvedMissing) != 1 || drift.ResolvedMissing[0] != "DB_HOST" {
		t.Errorf("expected ResolvedMissing=[DB_HOST], got %v", drift.ResolvedMissing)
	}
}

func TestCompareDrift_NewMismatched(t *testing.T) {
	old := makeSnap(nil, nil, nil)
	new := makeSnap(nil, nil, []diff.Mismatch{{Key: "PORT", LeftValue: "8080", RightValue: "9090"}})
	drift := snapshot.CompareDrift(old, new)
	if len(drift.NewMismatched) != 1 || drift.NewMismatched[0] != "PORT" {
		t.Errorf("expected NewMismatched=[PORT], got %v", drift.NewMismatched)
	}
}

func TestCompareDrift_IsClean_False(t *testing.T) {
	old := makeSnap(nil, nil, nil)
	new := makeSnap([]string{"SECRET"}, nil, nil)
	drift := snapshot.CompareDrift(old, new)
	if drift.IsClean() {
		t.Error("expected drift to not be clean")
	}
}
