package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/cronparse/internal/expr"
	"github.com/yourorg/cronparse/internal/human"
	"github.com/yourorg/cronparse/internal/next"
)

func main() {
	nFlag := flag.Int("n", 5, "number of next run times to display")
	fromFlag := flag.String("from", "", "start time in RFC3339 format (default: now)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cronparse [flags] <cron expression>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  cronparse \"*/5 9-17 * * 1-5\"\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	raw := flag.Arg(0)
	e, err := expr.Parse(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid cron expression: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Expression : %s\n", raw)
	fmt.Printf("Description: %s\n\n", human.Describe(e))

	from := time.Now()
	if *fromFlag != "" {
		from, err = time.Parse(time.RFC3339, *fromFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid --from time: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Next %d run(s) after %s:\n", *nFlag, from.Format(time.RFC3339))
	t := from
	for i := 0; i < *nFlag; i++ {
		t = next.After(e, t)
		fmt.Printf("  %d. %s\n", i+1, t.Format("2006-01-02 15:04 (Mon)"))
	}
}
