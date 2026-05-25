// Package schedule provides multi-expression scheduling support,
// allowing a named set of cron expressions to be evaluated together.
package schedule

import (
	"fmt"
	"sort"
	"time"

	"github.com/yourorg/cronparse/internal/expr"
	"github.com/yourorg/cronparse/internal/next"
)

// Entry pairs a human-readable name with a parsed cron expression.
type Entry struct {
	Name string
	Raw  string
	expr *expr.Expr
}

// Schedule holds a collection of named cron entries.
type Schedule struct {
	entries []*Entry
}

// New creates an empty Schedule.
func New() *Schedule {
	return &Schedule{}
}

// Add parses raw and registers it under name.
func (s *Schedule) Add(name, raw string) error {
	e, err := expr.Parse(raw)
	if err != nil {
		return fmt.Errorf("schedule %q: %w", name, err)
	}
	s.entries = append(s.entries, &Entry{Name: name, Raw: raw, expr: e})
	return nil
}

// NextAll returns each entry paired with its next run time after t,
// sorted ascending by next-run time.
func (s *Schedule) NextAll(t time.Time) []NextEntry {
	results := make([]NextEntry, 0, len(s.entries))
	for _, e := range s.entries {
		nxt := next.After(e.expr, t)
		results = append(results, NextEntry{Entry: e, Next: nxt})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Next.Before(results[j].Next)
	})
	return results
}

// NextEntry associates an Entry with its upcoming run time.
type NextEntry struct {
	*Entry
	Next time.Time
}

// Len returns the number of registered entries.
func (s *Schedule) Len() int { return len(s.entries) }
