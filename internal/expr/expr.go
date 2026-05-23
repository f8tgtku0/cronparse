// Package expr parses a full five-field cron expression into structured fields.
package expr

import (
	"fmt"
	"strings"

	"github.com/cronparse/internal/field"
)

// Expression holds the five parsed fields of a standard cron expression.
type Expression struct {
	Minute     *field.Field
	Hour       *field.Field
	DayOfMonth *field.Field
	Month      *field.Field
	DayOfWeek  *field.Field
	Raw        string
}

// Parse parses a standard 5-field cron expression string.
// Format: "<minute> <hour> <day-of-month> <month> <day-of-week>"
func Parse(expr string) (*Expression, error) {
	expr = strings.TrimSpace(expr)
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	kinds := []field.Kind{
		field.Minute,
		field.Hour,
		field.DayOfMonth,
		field.Month,
		field.DayOfWeek,
	}

	fields := make([]*field.Field, 5)
	for i, kind := range kinds {
		f, err := field.Parse(parts[i], kind)
		if err != nil {
			return nil, fmt.Errorf("parse error in field %d (%q): %w", i+1, parts[i], err)
		}
		fields[i] = f
	}

	return &Expression{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Raw:        expr,
	}, nil
}

// String returns a human-readable summary of the expression.
func (e *Expression) String() string {
	return fmt.Sprintf(
		"Cron(%s) => minute=%v hour=%v dom=%v month=%v dow=%v",
		e.Raw,
		e.Minute.Values,
		e.Hour.Values,
		e.DayOfMonth.Values,
		e.Month.Values,
		e.DayOfWeek.Values,
	)
}
