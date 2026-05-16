// Package profiler provides statistical analysis for parsed env maps.
//
// It is useful for quickly understanding the shape of an environment
// configuration: how many keys exist, how many are empty, which keys
// are likely sensitive, and which naming prefixes dominate.
//
// Basic usage:
//
//	env := map[string]string{
//		"DB_HOST":     "localhost",
//		"DB_PASSWORD": "s3cr3t",
//		"APP_NAME":    "myapp",
//	}
//
//	p := profiler.Analyze(env)
//	fmt.Println(p.TotalKeys)     // 3
//	fmt.Println(p.SensitiveKeys) // 1
//
//	top := profiler.TopPrefixes(p, 5)
//	fmt.Println(top) // [DB APP]
//
// The package relies on [redactor.IsSensitiveKey] for sensitivity
// detection, so custom patterns configured in the redactor will be
// reflected here automatically.
package profiler
