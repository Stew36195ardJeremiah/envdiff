package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	leftName := flag.String("left-name", "left", "label for the left env file")
	rightName := flag.String("right-name", "right", "label for the right env file")
	quiet := flag.Bool("quiet", false, "only print summary line")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <left.env> <right.env>\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[0], err)
		os.Exit(1)
	}

	right, err := parser.ParseFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[1], err)
		os.Exit(1)
	}

	result := diff.Compare(left, right)

	if *quiet {
		fmt.Println(reporter.Summary(result))
	} else {
		reporter.Report(os.Stdout, result, *leftName, *rightName)
	}

	if !result.IsClean() {
		os.Exit(1)
	}
}
