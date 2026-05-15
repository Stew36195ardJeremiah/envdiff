// Package merger provides functionality to merge multiple env maps
// into a single map, with configurable conflict resolution strategies.
package merger

// Strategy defines how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first map that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last map that defines the key.
	StrategyLast
)

// Options controls merge behaviour.
type Options struct {
	Strategy Strategy
	// TrackOrigin records which source index each key came from.
	TrackOrigin bool
}

// Result holds the merged env map and optional origin tracking.
type Result struct {
	Env    map[string]string
	// Origin maps each key to the index of the source map it came from.
	// Only populated when Options.TrackOrigin is true.
	Origin map[string]int
}

// Merge combines multiple env maps according to the given options.
// Sources are processed left to right.
func Merge(sources []map[string]string, opts Options) Result {
	env := make(map[string]string)
	origin := make(map[string]int)

	for i, src := range sources {
		for k, v := range src {
			_, exists := env[k]
			switch {
			case !exists:
				env[k] = v
				origin[k] = i
			case opts.Strategy == StrategyLast:
				env[k] = v
				origin[k] = i
			// StrategyFirst: keep existing value, do nothing
			}
		}
	}

	result := Result{Env: env}
	if opts.TrackOrigin {
		result.Origin = origin
	}
	return result
}
