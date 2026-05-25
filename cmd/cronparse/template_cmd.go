package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/yourorg/cronparse/internal/template"
)

// runTemplate handles the `template` subcommand, which lists available named
// cron expression templates or looks up a specific one by name.
//
// Usage:
//
//	cronparse template              # list all templates
//	cronparse template <name>       # show a specific template
func runTemplate(args []string) int {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cronparse template [name]")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "List all built-in cron expression templates, or look up one by name.")
		fmt.Fprintln(os.Stderr)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}

	switch fs.NArg() {
	case 0:
		return runTemplateList()
	case 1:
		return runTemplateLookup(fs.Arg(0))
	default:
		fmt.Fprintf(os.Stderr, "error: too many arguments (expected 0 or 1, got %d)\n", fs.NArg())
		fs.Usage()
		return 2
	}
}

// runTemplateList prints all available templates in a formatted table.
func runTemplateList() int {
	names := template.Names()
	if len(names) == 0 {
		fmt.Println("No templates available.")
		return 0
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tEXPRESSION\tDESCRIPTION")
	fmt.Fprintln(w, strings.Repeat("-", 12)+"\t"+strings.Repeat("-", 12)+"\t"+strings.Repeat("-", 30))

	for _, name := range names {
		tmpl, ok := template.Lookup(name)
		if !ok {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", name, tmpl.Expression, tmpl.Description)
	}

	w.Flush()
	return 0
}

// runTemplateLookup prints details for a single named template.
func runTemplateLookup(name string) int {
	tmpl, ok := template.Lookup(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown template %q\n", name)
		fmt.Fprintln(os.Stderr, "Run `cronparse template` to see available templates.")
		return 1
	}

	fmt.Printf("Name:        %s\n", tmpl.Name)
	fmt.Printf("Expression:  %s\n", tmpl.Expression)
	fmt.Printf("Description: %s\n", tmpl.Description)
	if tmpl.Notes != "" {
		fmt.Printf("Notes:       %s\n", tmpl.Notes)
	}
	return 0
}
