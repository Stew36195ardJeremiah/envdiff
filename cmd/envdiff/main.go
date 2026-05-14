package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/reporter"
)

func main() {
	var (
		prefix      = flag.String("prefix", "", "Only compare keys with this prefix")
		exclude     = flag.String("exclude", "", "Comma-separated list of keys to exclude")
		ignoreCase  = flag.Bool("ignore-case", false, "Case-insensitive key matching")
		maskSecrets = flag.Bool("mask", true, "Mask secret values in output")
		summaryOnly = flag.Bool("summary", false, "Print summary line only")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <left.env> <right.env>\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	leftEnv, err := parser.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[0], err)
		os.Exit(1)
	}

	rightEnv, err := parser.ParseFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[1], err)
		os.Exit(1)
	}

	var excludeKeys []string
	if *exclude != "" {
		for _, k := range strings.Split(*exclude, ",") {
			if k = strings.TrimSpace(k); k != "" {
				excludeKeys = append(excludeKeys, k)
			}
		}
	}

	opts := filter.Options{
		Prefix:      *prefix,
		ExcludeKeys: excludeKeys,
		IgnoreCase:  *ignoreCase,
	}

	leftEnv = filter.Apply(leftEnv, opts)
	rightEnv = filter.Apply(rightEnv, opts)

	result := diff.Compare(leftEnv, rightEnv)

	if *summaryOnly {
		fmt.Println(reporter.Summary(result))
	} else {
		reporter.Report(result, args[0], args[1], *maskSecrets, os.Stdout)
	}

	if !result.Clean() {
		os.Exit(1)
	}
}
