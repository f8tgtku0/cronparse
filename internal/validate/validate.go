// Package validate provides semantic validation for cron expressions beyond
// basic parsing — catching logical contradictions, unreachable schedules, and
// common misconfiguration patterns.
package validate

import (
	"fmt"
	"strings"

	"github.com/your-org/cronparse/internal/expr"
)

// Severity indicates how serious a validation finding is.
type Severity int

const (
	SeverityWarning Severity = iota
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}

// Finding represents a single validation result.
type Finding struct {
	Severity Severity
	Field    string
	Message  string
}

func (f Finding) String() string {
	return fmt.Sprintf("%s [%s]: %s", f.Severity, f.Field, f.Message)
}

// Result holds all findings for a validated expression.
type Result struct {
	Expression string
	Findings   []Finding
}

// OK returns true when there are no error-level findings.
func (r Result) OK() bool {
	for _, f := range r.Findings {
		if f.Severity == SeverityError {
			return false
		}
	}
	return true
}

// Summary returns a human-readable summary line.
func (r Result) Summary() string {
	if len(r.Findings) == 0 {
		return "expression is valid"
	}
	var parts []string
	for _, f := range r.Findings {
		parts = append(parts, f.String())
	}
	return strings.Join(parts, "; ")
}

// Validate parses and semantically validates a standard 5-field cron expression.
// It returns a Result containing any warnings or errors found.
func Validate(expression string) Result {
	r := Result{Expression: expression}

	e, err := expr.Parse(expression)
	if err != nil {
		r.Findings = append(r.Findings, Finding{
			Severity: SeverityError,
			Field:    "expression",
			Message:  fmt.Sprintf("parse error: %v", err),
		})
		return r
	}

	// Warn when both day-of-month and day-of-week are restricted simultaneously;
	// POSIX cron treats this as a union which is often unintentional.
	minField := "*"
	hourField := "*"
	domField := "*"
	dowField := "*"

	fields := strings.Fields(expression)
	if len(fields) == 5 {
		minField = fields[0]
		hourField = fields[1]
		domField = fields[2]
		dowField = fields[4]
	}
	_ = minField
	_ = hourField

	if domField != "*" && domField != "?" && dowField != "*" && dowField != "?" {
		r.Findings = append(r.Findings, Finding{
			Severity: SeverityWarning,
			Field:    "day-of-month / day-of-week",
			Message:  "both day-of-month and day-of-week are restricted; POSIX cron unions them, which may produce unexpected run times",
		})
	}

	// Warn on February 30/31 — will never match.
	if hasValue(e.Month, 2) && (hasValue(e.Dom, 30) || hasValue(e.Dom, 31)) {
		r.Findings = append(r.Findings, Finding{
			Severity: SeverityWarning,
			Field:    "day-of-month",
			Message:  "day 30 or 31 in February will never occur",
		})
	}

	// Warn on 31st day in months that have only 30 days.
	shortMonths := []int{4, 6, 9, 11} // April, June, September, November
	for _, m := range shortMonths {
		if hasValue(e.Month, m) && hasValue(e.Dom, 31) {
			r.Findings = append(r.Findings, Finding{
				Severity: SeverityWarning,
				Field:    "day-of-month",
				Message:  fmt.Sprintf("day 31 never occurs in month %d", m),
			})
			break
		}
	}

	return r
}

// hasValue reports whether the sorted int slice contains v.
func hasValue(values []int, v int) bool {
	for _, x := range values {
		if x == v {
			return true
		}
	}
	return false
}
