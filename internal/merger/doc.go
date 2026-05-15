// Package merger combines multiple parsed env maps into a single unified map.
//
// It is useful when an application's configuration is spread across several
// .env files (e.g. a base .env, a .env.local override, and a secrets file)
// and you want to reason about the effective configuration as a whole.
//
// # Strategies
//
// Two conflict-resolution strategies are supported:
//
//   - StrategyFirst (default): the value from the earliest source wins.
//     This mirrors how many tools (e.g. docker-compose) handle overrides.
//
//   - StrategyLast: the value from the latest source wins.
//     Useful when later files are intended to override earlier ones.
//
// # Origin Tracking
//
// When Options.TrackOrigin is enabled, the Result.Origin map records the
// zero-based index of the source slice that contributed each key. This is
// helpful for debugging unexpected value origins in deep override chains.
//
// # Example
//
//	result := merger.Merge([]map[string]string{base, local}, merger.Options{
//		Strategy:    merger.StrategyLast,
//		TrackOrigin: true,
//	})
package merger
