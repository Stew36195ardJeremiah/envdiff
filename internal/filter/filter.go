package filter

import (
	"strings"
)

// Options holds filtering configuration.
type Options struct {
	// Prefix restricts comparison to keys with this prefix.
	Prefix string
	// ExcludeKeys is a set of exact key names to ignore.
	ExcludeKeys []string
	// IgnoreCase treats key names case-insensitively.
	IgnoreCase bool
}

// Apply returns a filtered copy of the input map based on the given Options.
func Apply(env map[string]string, opts Options) map[string]string {
	excluded := make(map[string]string, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		norm := normalizeKey(k, opts.IgnoreCase)
		excluded[norm] = ""
	}

	out := make(map[string]string)
	for k, v := range env {
		norm := normalizeKey(k, opts.IgnoreCase)

		if opts.Prefix != "" {
			pfx := normalizeKey(opts.Prefix, opts.IgnoreCase)
			if !strings.HasPrefix(norm, pfx) {
				continue
			}
		}

		if _, skip := excluded[norm]; skip {
			continue
		}

		out[k] = v
	}
	return out
}

// MatchesPrefix reports whether the key matches the given prefix,
// respecting the IgnoreCase flag.
func MatchesPrefix(key, prefix string, ignoreCase bool) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(
		normalizeKey(key, ignoreCase),
		normalizeKey(prefix, ignoreCase),
	)
}

func normalizeKey(k string, ignoreCase bool) string {
	if ignoreCase {
		return strings.ToLower(k)
	}
	return k
}
