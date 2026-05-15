// Package snapshot provides functionality to save and load env diff results
// to/from disk, enabling comparison of results over time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot represents a saved diff result with metadata.
type Snapshot struct {
	CreatedAt time.Time  `json:"created_at"`
	Label     string     `json:"label,omitempty"`
	Result    diff.Result `json:"result"`
}

// Save writes a snapshot of the given diff result to the specified file path.
func Save(path, label string, result diff.Result) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Result:    result,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return snap, nil
}
