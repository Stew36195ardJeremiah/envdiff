// Package snapshot provides save/load functionality for envdiff results,
// allowing teams to persist a diff result to disk and later compare it
// against a newer result to detect drift over time.
//
// # Saving a Snapshot
//
//	err := snapshot.Save("baseline.json", "v1.0-deploy", diffResult)
//
// # Loading a Snapshot
//
//	snap, err := snapshot.Load("baseline.json")
//
// # Detecting Drift
//
// Use CompareDrift to identify keys that have become missing or mismatched
// since the baseline snapshot was taken:
//
//	drift := snapshot.CompareDrift(olderSnap, newerSnap)
//	if !drift.IsClean() {
//		fmt.Println("New issues detected:", drift.NewMissing)
//	}
package snapshot
