// Package history provides utilities for computing previous run times
// of a cron expression relative to a given reference time.
package history

import (
	"fmt"
	"time"

	"github.com/your-org/cronparse/internal/expr"
	"github.com/your-org/cronparse/internal/next"
)

// Before returns the n most recent times the cron expression would have
// fired strictly before t, ordered from most recent to oldest.
// Returns an error if the expression is invalid or n < 1.
func Before(expression string, t time.Time, n int) ([]time.Time, error) {
	if n < 1 {
		return nil, fmt.Errorf("history: n must be at least 1, got %d", n)
	}

	e, err := expr.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("history: %w", err)
	}

	results := make([]time.Time, 0, n)

	// Walk backwards in one-minute steps from t-1m, collecting matches.
	// We search up to 4 years back to avoid an infinite loop on
	// pathological expressions.
	cursor := t.Add(-time.Minute).Truncate(time.Minute)
	limit := t.Add(-4 * 365 * 24 * time.Hour)

	for len(results) < n && cursor.After(limit) {
		// next.After finds the next run at or after cursor; if it equals
		// cursor then cursor itself is a valid firing time.
		candidate := next.After(e, cursor)
		if candidate.Equal(cursor) {
			results = append(results, cursor)
		}
		cursor = cursor.Add(-time.Minute)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("history: no occurrences found in the last 4 years")
	}

	return results, nil
}
