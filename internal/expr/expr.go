package expr

import (
	"fmt"
	"strings"

	"github.com/cronparse/internal/field"
)

// Expr holds the five parsed cron fields.
type Expr struct {
	Minute     *field.Field
	Hour       *field.Field
	DayOfMonth *field.Field
	Month      *field.Field
	DayOfWeek  *field.Field
}

// fieldSpec describes the valid bounds for each cron position.
type fieldSpec struct {
	name string
	min  int
	max  int
}

var specs = []fieldSpec{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"day-of-month", 1, 31},
	{"month", 1, 12},
	{"day-of-week", 0, 6},
}

// Parse parses a standard five-field cron expression string.
func Parse(s string) (*Expr, error) {
	parts := strings.Fields(s)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expr: expected 5 fields, got %d", len(parts))
	}

	fields := make([]*field.Field, 5)
	for i, p := range parts {
		f, err := field.Parse(p, specs[i].min, specs[i].max)
		if err != nil {
			return nil, fmt.Errorf("expr: %s field: %w", specs[i].name, err)
		}
		fields[i] = f
	}

	return &Expr{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
	}, nil
}

// String reassembles the expression in canonical form.
func (e *Expr) String() string {
	return strings.Join([]string{
		e.Minute.String(),
		e.Hour.String(),
		e.DayOfMonth.String(),
		e.Month.String(),
		e.DayOfWeek.String(),
	}, " ")
}
