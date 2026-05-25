// Package convert translates cron expressions between different dialects
// (standard 5-field, Quartz 6-field with seconds, AWS EventBridge, etc.).
package convert

import (
	"fmt"
	"strings"
)

// Dialect represents a supported cron dialect.
type Dialect string

const (
	DialectStandard   Dialect = "standard"   // min hour dom month dow
	DialectQuartz     Dialect = "quartz"     // sec min hour dom month dow [year]
	DialectEventBridge Dialect = "eventbridge" // min hour dom month dow year
)

// Result holds the converted expression and any informational notes.
type Result struct {
	Expression string
	Notes      []string
}

// Convert translates expr from src dialect to dst dialect.
func Convert(expr string, src, dst Dialect) (*Result, error) {
	if src == dst {
		return &Result{Expression: expr}, nil
	}

	fields, err := splitFields(expr)
	if err != nil {
		return nil, err
	}

	// Normalise to standard 5-field first.
	std, notes, err := toStandard(fields, src)
	if err != nil {
		return nil, err
	}

	// Convert from standard to target dialect.
	out, extraNotes, err := fromStandard(std, dst)
	if err != nil {
		return nil, err
	}

	return &Result{
		Expression: out,
		Notes:      append(notes, extraNotes...),
	}, nil
}

func splitFields(expr string) ([]string, error) {
	f := strings.Fields(expr)
	if len(f) < 5 || len(f) > 7 {
		return nil, fmt.Errorf("convert: unexpected field count %d in %q", len(f), expr)
	}
	return f, nil
}

// toStandard returns a 5-element slice [min hour dom month dow].
func toStandard(fields []string, src Dialect) ([]string, []string, error) {
	var notes []string
	switch src {
	case DialectStandard:
		if len(fields) != 5 {
			return nil, nil, fmt.Errorf("convert: standard dialect requires 5 fields, got %d", len(fields))
		}
		return fields, notes, nil
	case DialectQuartz:
		if len(fields) < 6 {
			return nil, nil, fmt.Errorf("convert: quartz dialect requires at least 6 fields, got %d", len(fields))
		}
		if fields[0] != "0" && fields[0] != "*" {
			notes = append(notes, fmt.Sprintf("seconds field %q dropped; standard cron has no seconds", fields[0]))
		}
		// drop seconds (index 0) and optional year (last field if 7)
		return fields[1:6], notes, nil
	case DialectEventBridge:
		if len(fields) != 6 {
			return nil, nil, fmt.Errorf("convert: eventbridge dialect requires 6 fields, got %d", len(fields))
		}
		if fields[5] != "*" {
			notes = append(notes, fmt.Sprintf("year field %q dropped; standard cron has no year", fields[5]))
		}
		return fields[:5], notes, nil
	}
	return nil, nil, fmt.Errorf("convert: unknown source dialect %q", src)
}

// fromStandard builds the target dialect string from a 5-field standard slice.
func fromStandard(std []string, dst Dialect) (string, []string, error) {
	var notes []string
	switch dst {
	case DialectStandard:
		return strings.Join(std, " "), notes, nil
	case DialectQuartz:
		// Prepend seconds=0, append year=*
		fields := append([]string{"0"}, std...)
		fields = append(fields, "*")
		notes = append(notes, "seconds field set to 0; year field set to *")
		return strings.Join(fields, " "), notes, nil
	case DialectEventBridge:
		fields := append(std, "*")
		notes = append(notes, "year field set to *")
		return strings.Join(fields, " "), notes, nil
	}
	return "", nil, fmt.Errorf("convert: unknown destination dialect %q", dst)
}
