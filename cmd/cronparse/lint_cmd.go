package main

import (
	"fmt"
	"os"

	"github.com/yourorg/cronparse/internal/lint"
)

// runLint checks the given expression and prints any diagnostics to stdout.
// It exits with a non-zero status when at least one Error diagnostic is found.
func runLint(expression string) {
	diags := lint.Check(expression)

	if len(diags) == 0 {
		fmt.Println("OK: no issues found")
		return
	}

	hasError := false
	for _, d := range diags {
		fmt.Println(d)
		if d.Severity == lint.Error {
			hasError = true
		}
	}

	if hasError {
		os.Exit(2)
	}
}
