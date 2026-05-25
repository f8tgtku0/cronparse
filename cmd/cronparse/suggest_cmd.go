package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/cronparse/internal/suggest"
)

// runSuggest handles the "suggest" subcommand, which searches for cron
// expressions matching a natural-language query or lists all built-in
// suggestions when no query is provided.
//
// Usage:
//
//	cronparse suggest [flags] [query]
//
// Flags:
//
//	-n int   maximum number of results to display (default 5)
func runSuggest(args []string) int {
	fs := flag.NewFlagSet("suggest", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	maxResults := fs.Int("n", 5, "maximum number of results to display")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	query := strings.TrimSpace(strings.Join(fs.Args(), " "))

	var results []suggest.Entry

	if query == "" {
		// No query — list all known suggestions.
		results = suggest.All()
	} else {
		// Search by keyword match first; fall back to closest by edit distance.
		results = suggest.Search(query)
		if len(results) == 0 {
			results = suggest.Closest(query, *maxResults)
			if len(results) > 0 {
				fmt.Fprintf(os.Stderr, "No exact matches for %q. Did you mean:\n", query)
			}
		}
	}

	if len(results) == 0 {
		fmt.Fprintf(os.Stderr, "no suggestions found for %q\n", query)
		return 1
	}

	// Trim to requested maximum.
	if len(results) > *maxResults {
		results = results[:*maxResults]
	}

	// Print results in a table-like format.
	const exprWidth = 20
	fmt.Printf("%-*s  %s\n", exprWidth, "EXPRESSION", "DESCRIPTION")
	fmt.Printf("%-*s  %s\n", exprWidth, strings.Repeat("-", exprWidth), strings.Repeat("-", 40))
	for _, e := range results {
		fmt.Printf("%-*s  %s\n", exprWidth, e.Expression, e.Description)
	}

	return 0
}

func init() {
	registerCommand("suggest", "search built-in cron expression suggestions", runSuggest)
}
