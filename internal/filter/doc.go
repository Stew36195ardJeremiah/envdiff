// Package filter provides utilities for narrowing the set of environment
// variables considered during a diff operation.
//
// It supports three independent filtering mechanisms that can be combined:
//
//   - Prefix filtering: only keys whose names start with a given prefix are
//     retained. Useful when a single .env file contains variables for multiple
//     services (e.g. "APP_", "DB_", "CACHE_").
//
//   - Key exclusion: an explicit list of key names that should be ignored
//     regardless of their values. Commonly used to suppress known secrets such
//     as DB_PASSWORD or SECRET_TOKEN from the diff output.
//
//   - Case-insensitive matching: when enabled, both prefix and exclusion
//     comparisons are performed after lowercasing the key names, allowing
//     filters written in any case to match keys regardless of their original
//     casing in the file.
//
// Typical usage:
//
//	opts := filter.Options{
//		Prefix:      "APP_",
//		ExcludeKeys: []string{"APP_SECRET"},
//	}
//	filtered := filter.Apply(rawEnv, opts)
package filter
