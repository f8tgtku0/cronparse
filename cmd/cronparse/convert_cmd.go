package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/cronparse/internal/convert"
)

// runConvert implements the `cronparse convert` sub-command.
// Usage: cronparse convert -from <dialect> -to <dialect> <expression>
func runConvert(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	fs.SetOutput(out)

	fromFlag := fs.String("from", "standard", "source dialect (standard|quartz|eventbridge)")
	toFlag := fs.String("to", "standard", "target dialect (standard|quartz|eventbridge)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(out, "error: expression argument required")
		fmt.Fprintln(out, "usage: cronparse convert -from <dialect> -to <dialect> <expression>")
		return 2
	}

	expr := fs.Arg(0)
	src := convert.Dialect(*fromFlag)
	dst := convert.Dialect(*toFlag)

	res, err := convert.Convert(expr, src, dst)
	if err != nil {
		fmt.Fprintf(out, "error: %v\n", err)
		return 1
	}

	fmt.Fprintf(out, "Result: %s\n", res.Expression)
	for _, note := range res.Notes {
		fmt.Fprintf(out, "  note: %s\n", note)
	}
	return 0
}

// init registers the convert sub-command when the binary starts.
func init() {
	_ = os.Args // referenced to satisfy import; registration happens in main.
}
