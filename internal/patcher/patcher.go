// Package patcher applies a set of key-value updates to an existing env map,
// supporting upsert (add or overwrite) and delete operations.
package patcher

import "fmt"

// Operation defines the type of patch action.
type Operation string

const (
	OpUpsert Operation = "upsert"
	OpDelete Operation = "delete"
)

// Patch represents a single change to apply to an env map.
type Patch struct {
	Key   string
	Value string
	Op    Operation
}

// Result holds the outcome of applying patches.
type Result struct {
	Env     map[string]string
	Applied []Patch
	Skipped []Patch
}

// Options controls patcher behaviour.
type Options struct {
	// FailOnMissing returns an error if a delete patch targets a key that does not exist.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
	}
}

// Apply applies the given patches to a copy of env and returns the result.
func Apply(env map[string]string, patches []Patch, opts Options) (Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := Result{Env: out}

	for _, p := range patches {
		switch p.Op {
		case OpUpsert:
			out[p.Key] = p.Value
			result.Applied = append(result.Applied, p)

		case OpDelete:
			if _, exists := out[p.Key]; !exists {
				if opts.FailOnMissing {
					return Result{}, fmt.Errorf("patcher: delete target %q not found", p.Key)
				}
				result.Skipped = append(result.Skipped, p)
				continue
			}
			delete(out, p.Key)
			result.Applied = append(result.Applied, p)

		default:
			return Result{}, fmt.Errorf("patcher: unknown operation %q for key %q", p.Op, p.Key)
		}
	}

	return result, nil
}
