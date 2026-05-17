// Package annotator attaches inline comments to env file entries,
// describing their diff status (missing, mismatched, ok).
package annotator

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Status represents the annotation state of a key.
type Status string

const (
	StatusOK         Status = "ok"
	StatusMissing    Status = "missing"
	StatusMismatched Status = "mismatched"
)

// Annotation holds a key, its value, and its diff status.
type Annotation struct {
	Key     string
	Value   string
	Status  Status
	Comment string
}

// Annotate takes a map of env values and a diff.Result, and returns
// a slice of Annotations sorted by key.
func Annotate(env map[string]string, result diff.Result) []Annotation {
	missingSet := toSet(result.MissingInRight)
	mismatchedSet := make(map[string]diff.Mismatch)
	for _, m := range result.Mismatched {
		mismatchedSet[m.Key] = m
	}

	annotations := make([]Annotation, 0, len(env))
	for k, v := range env {
		a := Annotation{Key: k, Value: v}
		switch {
		case missingSet[k]:
			a.Status = StatusMissing
			a.Comment = "# missing in right"
		case mismatchedSet[k].Key != "":
			a.Status = StatusMismatched
			a.Comment = fmt.Sprintf("# mismatched: right=%q", mismatchedSet[k].RightVal)
		default:
			a.Status = StatusOK
			a.Comment = "# ok"
		}
		annotations = append(annotations, a)
	}

	sort.Slice(annotations, func(i, j int) bool {
		return annotations[i].Key < annotations[j].Key
	})
	return annotations
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
