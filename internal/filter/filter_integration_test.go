package filter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/parser"
)

func TestFilter_WithParsedFile(t *testing.T) {
	env, err := parser.ParseFile("../parser/testdata/sample.env")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	if len(env) == 0 {
		t.Skip("sample.env has no entries, skipping integration test")
	}

	// Apply no filters — should return all keys.
	all := filter.Apply(env, filter.Options{})
	if len(all) != len(env) {
		t.Errorf("expected %d keys unfiltered, got %d", len(env), len(all))
	}

	// Exclude every key explicitly — result should be empty.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	none := filter.Apply(env, filter.Options{ExcludeKeys: keys})
	if len(none) != 0 {
		t.Errorf("expected 0 keys after excluding all, got %d", len(none))
	}
}
