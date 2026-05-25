package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/cronparse/internal/diff"
)

// runDiff implements the "diff" sub-command. It expects exactly two positional
// arguments: the two cron expressions to compare.
//
// Usage:
//
//	cronparse diff "0 9 * * 1" "30 17 * * 5"
func runDiff(args []string) int {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cronparse diff <expr-a> <expr-b>")
		fmt.Fprintln(os.Stderr, "  Compare two cron expressions and describe what changed.")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fs.NArg() != 2 {
		fs.Usage()
		return 2
	}

	rawA := fs.Arg(0)
	rawB := fs.Arg(1)

	res, err := diff.Compare(rawA, rawB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println(res.Summary)
	if len(res.Fields) > 0 {
		fmt.Println()
		fmt.Println("Changed fields:")
		for _, f := range res.Fields {
			fmt.Printf("  %-16s %s  ->  %s\n", f.Name+":", f.From, f.To)
		}
	}
	return 0
}
