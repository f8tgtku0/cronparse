// Package lint provides validation and warning diagnostics for cron expressions.
package lint

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/internal/expr"
)

// Severity indicates how serious a diagnostic is.
type Severity int

const (
	Warning Severity = iota
	Error
)

func (s Severity) String() string {
	switch s {
	case Warning:
		return "warning"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Diagnostic holds a single lint result.
type Diagnostic struct {
	Severity Severity
	Message  string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: %s", d.Severity, d.Message)
}

// Check parses the given cron expression and returns any diagnostics.
// A parse failure produces an Error diagnostic. Suspicious but valid
// expressions produce Warning diagnostics.
func Check(expression string) []Diagnostic {
	var diags []Diagnostic

	e, err := expr.Parse(expression)
	if err != nil {
		diags = append(diags, Diagnostic{Error, err.Error()})
		return diags
	}

	// Warn when every field is a wildcard — likely a placeholder.
	fields := strings.Fields(expression)
	allWild := true
	for _, f := range fields {
		if f != "*" {
			allWild = false
			break
		}
	}
	if allWild {
		diags = append(diags, Diagnostic{Warning, "expression matches every minute — is this intentional?"})
	}

	// Warn when both day-of-month and day-of-week are constrained.
	_ = e
	if len(fields) == 5 && fields[2] != "*" && fields[4] != "*" {
		diags = append(diags, Diagnostic{Warning, "both day-of-month and day-of-week are set; most cron implementations use OR semantics"})
	}

	return diags
}
