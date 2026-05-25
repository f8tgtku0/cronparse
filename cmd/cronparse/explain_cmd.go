package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/yourorg/cronparse/internal/explain"
)

func runExplain(args []string) error {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	jsonOut := fs.Bool("json", false, "output result as JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("explain: usage: cronparse explain [--json] '<expression>'")
	}

	expr := fs.Arg(0)
	res, err := explain.Explain(expr)
	if err != nil {
		return err
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FIELD\tRAW\tMEANING")
	fmt.Fprintln(w, "-----\t---\t-------")
	for _, f := range res.Fields {
		fmt.Fprintf(w, "%s\t%s\t%s\n", f.Field, f.Raw, f.Meaning)
	}
	w.Flush()
	fmt.Printf("\nSummary: %s\n", res.Summary)
	return nil
}
