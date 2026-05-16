// Package profiler analyses an env map and produces a statistical
// profile: key count, empty-value count, sensitive-key count, and a
// breakdown of key prefixes found in the data.
package profiler

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/redactor"
)

// Profile holds summary statistics for a parsed env map.
type Profile struct {
	TotalKeys     int
	EmptyValues   int
	SensitiveKeys int
	// PrefixCounts maps each detected prefix (e.g. "DB", "AWS") to the
	// number of keys that share it.
	PrefixCounts map[string]int
}

// Analyze computes a Profile for the provided env map.
func Analyze(env map[string]string) Profile {
	p := Profile{
		TotalKeys:    len(env),
		PrefixCounts: make(map[string]int),
	}

	for k, v := range env {
		if strings.TrimSpace(v) == "" {
			p.EmptyValues++
		}
		if redactor.IsSensitiveKey(k) {
			p.SensitiveKeys++
		}
		if prefix := extractPrefix(k); prefix != "" {
			p.PrefixCounts[prefix]++
		}
	}

	return p
}

// TopPrefixes returns up to n prefixes ordered by descending count.
// Ties are broken alphabetically.
func TopPrefixes(p Profile, n int) []string {
	type kv struct {
		key   string
		count int
	}
	pairs := make([]kv, 0, len(p.PrefixCounts))
	for k, c := range p.PrefixCounts {
		pairs = append(pairs, kv{k, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	result := make([]string, 0, n)
	for i := 0; i < len(pairs) && i < n; i++ {
		result = append(result, pairs[i].key)
	}
	return result
}

// extractPrefix returns the portion of key before the first underscore,
// provided the key contains at least one underscore and the prefix is
// non-empty.
func extractPrefix(key string) string {
	idx := strings.Index(key, "_")
	if idx <= 0 {
		return ""
	}
	return key[:idx]
}
