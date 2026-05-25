// Package explain provides a field-by-field breakdown of a cron expression,
// returning structured annotations suitable for display or tooling.
package explain

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/internal/expr"
)

// FieldExplanation holds the name and human-readable description of a single
// cron field.
type FieldExplanation struct {
	Field   string `json:"field"`
	Raw     string `json:"raw"`
	Meaning string `json:"meaning"`
}

// Result is the full explanation of a parsed cron expression.
type Result struct {
	Expression string             `json:"expression"`
	Fields     []FieldExplanation `json:"fields"`
	Summary    string             `json:"summary"`
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Explain parses the given cron expression and returns a structured breakdown
// of each field. Returns an error if the expression is invalid.
func Explain(expression string) (*Result, error) {
	e, err := expr.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("explain: %w", err)
	}

	parts := strings.Fields(expression)
	if len(parts) != 5 {
		return nil, fmt.Errorf("explain: expected 5 fields, got %d", len(parts))
	}

	fields := make([]FieldExplanation, 5)
	for i, name := range fieldNames {
		fields[i] = FieldExplanation{
			Field:   name,
			Raw:     parts[i],
			Meaning: describeField(name, parts[i], e),
		}
	}

	return &Result{
		Expression: expression,
		Fields:     fields,
		Summary:    buildSummary(fields),
	}, nil
}

func describeField(name, raw string, _ *expr.Expr) string {
	if raw == "*" {
		return "every " + name
	}
	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		return fmt.Sprintf("every %s %s(s)", parts[1], name)
	}
	if strings.Contains(raw, "-") {
		parts := strings.SplitN(raw, "-", 2)
		return fmt.Sprintf("%s %s through %s", name, parts[0], parts[1])
	}
	if strings.Contains(raw, ",") {
		return fmt.Sprintf("%s in [%s]", name, raw)
	}
	return fmt.Sprintf("%s %s", name, raw)
}

func buildSummary(fields []FieldExplanation) string {
	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = f.Meaning
	}
	return strings.Join(parts, "; ")
}
