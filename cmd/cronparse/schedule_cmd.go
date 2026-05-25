package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/cronparse/internal/schedule"
)

// runSchedule handles the "schedule" subcommand, which manages a named
// collection of cron expressions stored in a JSON file. It supports
// adding entries, listing them, and computing the next N run times.
//
// Usage:
//
//	cronparse schedule -file <path> [-add <name>=<expr>] [-next <n>] [-list]
func runSchedule(args []string) error {
	fs := flag.NewFlagSet("schedule", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	filePath := fs.String("file", "schedule.json", "Path to the schedule JSON file")
	addEntry := fs.String("add", "", "Add a named entry in the form name=expression")
	listFlag := fs.Bool("list", false, "List all named entries in the schedule file")
	nextCount := fs.Int("next", 0, "Print the next N run times across all entries")
	fromStr := fs.String("from", "", "Reference time for -next (RFC3339); defaults to now")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load or create a fresh schedule.
	sched, err := schedule.LoadJSON(*filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading schedule: %w", err)
	}
	if sched == nil {
		sched = schedule.New()
	}

	// Add a new entry if requested.
	if *addEntry != "" {
		name, expr, ok := splitNameExpr(*addEntry)
		if !ok {
			return fmt.Errorf("-add value must be in the form name=expression, got: %q", *addEntry)
		}
		if err := sched.Add(name, expr); err != nil {
			return fmt.Errorf("adding entry %q: %w", name, err)
		}
		if err := schedule.SaveJSON(*filePath, sched); err != nil {
			return fmt.Errorf("saving schedule: %w", err)
		}
		fmt.Printf("Added entry %q to %s\n", name, *filePath)
	}

	// List all entries.
	if *listFlag {
		entries := sched.Entries()
		if len(entries) == 0 {
			fmt.Println("No entries in schedule.")
		} else {
			fmt.Printf("%-20s  %s\n", "NAME", "EXPRESSION")
			for _, e := range entries {
				fmt.Printf("%-20s  %s\n", e.Name, e.Raw)
			}
		}
	}

	// Print next N run times.
	if *nextCount > 0 {
		ref := time.Now()
		if *fromStr != "" {
			parsed, err := time.Parse(time.RFC3339, *fromStr)
			if err != nil {
				return fmt.Errorf("parsing -from time: %w", err)
			}
			ref = parsed
		}

		runs := sched.NextAll(ref, *nextCount)
		if len(runs) == 0 {
			fmt.Println("No upcoming runs (schedule may be empty).")
		} else {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(runs); err != nil {
				return fmt.Errorf("encoding output: %w", err)
			}
		}
	}

	return nil
}

// splitNameExpr splits a string of the form "name=expression" into its
// constituent parts. Returns ok=false if no '=' separator is found.
func splitNameExpr(s string) (name, expr string, ok bool) {
	for i, ch := range s {
		if ch == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
