// Package sorter provides utilities for sorting and grouping env key-value
// pairs by prefix, alphabetical order, or custom criteria.
package sorter

import (
	"sort"
	"strings"
)

// Options controls how sorting and grouping is applied.
type Options struct {
	// Alphabetical sorts keys A-Z when true (default: true).
	Alphabetical bool
	// GroupByPrefix groups keys sharing the same prefix (e.g. DB_, AWS_).
	GroupByPrefix bool
	// PrefixSeparator is the delimiter used to detect prefixes (default: "_").
	PrefixSeparator string
}

// DefaultOptions returns sensible defaults for sorting.
func DefaultOptions() Options {
	return Options{
		Alphabetical:    true,
		GroupByPrefix:   false,
		PrefixSeparator: "_",
	}
}

// Sort returns a new map with keys sorted according to opts. Because maps are
// unordered the function returns the sorted key slice alongside the original
// map so callers can iterate in order.
func Sort(env map[string]string, opts Options) (map[string]string, []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if opts.Alphabetical {
		sort.Strings(keys)
	}

	if opts.GroupByPrefix {
		keys = groupByPrefix(keys, opts.PrefixSeparator)
	}

	return env, keys
}

// GroupKeys partitions a sorted key slice into groups sharing the same prefix.
// Keys without a separator are placed in a group keyed by "".
func GroupKeys(keys []string, separator string) map[string][]string {
	groups := make(map[string][]string)
	for _, k := range keys {
		prefix := ""
		if idx := strings.Index(k, separator); idx > 0 {
			prefix = k[:idx]
		}
		groups[prefix] = append(groups[prefix], k)
	}
	return groups
}

// groupByPrefix reorders keys so that members of the same prefix appear
// consecutively, preserving alphabetical order within each group.
func groupByPrefix(keys []string, sep string) []string {
	groups := GroupKeys(keys, sep)

	// Collect and sort group prefixes for deterministic output.
	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	result := make([]string, 0, len(keys))
	for _, p := range prefixes {
		result = append(result, groups[p]...)
	}
	return result
}
