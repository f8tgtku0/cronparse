// Package diff computes the difference between two cron expressions,
// describing which fields changed and how the schedule shifts.
package diff

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/internal/expr"
	"github.com/yourorg/cronparse/internal/human"
)

// Result holds a human-readable summary of differences between two expressions.
type Result struct {
	Fields  []FieldDiff
	Summary string
}

// FieldDiff describes a change in a single cron field.
type FieldDiff struct {
	Name string
	From string
	To   string
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Compare returns a Result describing the differences between two raw cron
// expression strings. It returns an error if either expression is invalid.
func Compare(rawA, rawB string) (*Result, error) {
	a, err := expr.Parse(rawA)
	if err != nil {
		return nil, fmt.Errorf("first expression: %w", err)
	}
	b, err := expr.Parse(rawB)
	if err != nil {
		return nil, fmt.Errorf("second expression: %w", err)
	}

	partsA := strings.Fields(rawA)
	partsB := strings.Fields(rawB)
	_ = a
	_ = b

	var diffs []FieldDiff
	for i, name := range fieldNames {
		if partsA[i] != partsB[i] {
			diffs = append(diffs, FieldDiff{
				Name: name,
				From: partsA[i],
				To:   partsB[i],
			})
		}
	}

	summary := buildSummary(rawA, rawB, diffs)
	return &Result{Fields: diffs, Summary: summary}, nil
}

func buildSummary(rawA, rawB string, diffs []FieldDiff) string {
	if len(diffs) == 0 {
		return "Expressions are identical."
	}
	var sb strings.Builder
	descA := human.Describe(mustParseQuiet(rawA))
	descB := human.Describe(mustParseQuiet(rawB))
	fmt.Fprintf(&sb, "Schedule changed from %q to %q.\n", descA, descB)
	fmt.Fprintf(&sb, "%d field(s) differ: ", len(diffs))
	names := make([]string, len(diffs))
	for i, d := range diffs {
		names[i] = d.Name
	}
	sb.WriteString(strings.Join(names, ", ") + ".")
	return sb.String()
}

func mustParseQuiet(raw string) *expr.Expr {
	e, _ := expr.Parse(raw)
	return e
}
