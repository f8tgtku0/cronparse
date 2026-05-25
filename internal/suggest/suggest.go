// Package suggest provides cron expression suggestions based on natural
// language descriptions or common scheduling patterns.
package suggest

import (
	"fmt"
	"strings"
)

// Suggestion represents a suggested cron expression with a human-readable label.
type Suggestion struct {
	// Label is a short description of the schedule.
	Label string
	// Expression is the standard 5-field cron expression.
	Expression string
	// Notes contains optional clarifying information.
	Notes string
}

// String returns a formatted representation of the suggestion.
func (s Suggestion) String() string {
	if s.Notes != "" {
		return fmt.Sprintf("%-35s %s  # %s", s.Expression, s.Label, s.Notes)
	}
	return fmt.Sprintf("%-35s %s", s.Expression, s.Label)
}

// builtins is the catalog of well-known scheduling patterns.
var builtins = []Suggestion{
	{Label: "every minute", Expression: "* * * * *"},
	{Label: "every 5 minutes", Expression: "*/5 * * * *"},
	{Label: "every 10 minutes", Expression: "*/10 * * * *"},
	{Label: "every 15 minutes", Expression: "*/15 * * * *"},
	{Label: "every 30 minutes", Expression: "*/30 * * * *"},
	{Label: "every hour", Expression: "0 * * * *"},
	{Label: "every 2 hours", Expression: "0 */2 * * *"},
	{Label: "every 6 hours", Expression: "0 */6 * * *"},
	{Label: "every 12 hours", Expression: "0 */12 * * *"},
	{Label: "once a day at midnight", Expression: "0 0 * * *"},
	{Label: "once a day at noon", Expression: "0 12 * * *"},
	{Label: "every weekday at 9am", Expression: "0 9 * * 1-5"},
	{Label: "every weekday at midnight", Expression: "0 0 * * 1-5"},
	{Label: "every weekend at midnight", Expression: "0 0 * * 6,0"},
	{Label: "once a week on Monday", Expression: "0 0 * * 1"},
	{Label: "once a week on Sunday", Expression: "0 0 * * 0"},
	{Label: "first day of the month", Expression: "0 0 1 * *"},
	{Label: "last day of the month", Expression: "0 0 28-31 * *", Notes: "approximate; cron has no true last-day support"},
	{Label: "first day of every quarter", Expression: "0 0 1 1,4,7,10 *"},
	{Label: "once a year on Jan 1", Expression: "0 0 1 1 *"},
	{Label: "every hour during business hours", Expression: "0 9-17 * * 1-5"},
	{Label: "every 5 minutes during business hours", Expression: "*/5 9-17 * * 1-5"},
}

// All returns the full list of built-in suggestions.
func All() []Suggestion {
	result := make([]Suggestion, len(builtins))
	copy(result, builtins)
	return result
}

// Search returns suggestions whose label contains any of the provided keywords
// (case-insensitive). If no keywords are given, all suggestions are returned.
func Search(keywords ...string) []Suggestion {
	if len(keywords) == 0 {
		return All()
	}

	normalized := make([]string, len(keywords))
	for i, kw := range keywords {
		normalized[i] = strings.ToLower(strings.TrimSpace(kw))
	}

	var results []Suggestion
	for _, s := range builtins {
		label := strings.ToLower(s.Label)
		for _, kw := range normalized {
			if kw != "" && strings.Contains(label, kw) {
				results = append(results, s)
				break
			}
		}
	}
	return results
}

// Closest returns up to n suggestions that best match the query string by
// counting how many whitespace-separated tokens from the query appear in the
// suggestion label. Ties are broken by original catalog order.
func Closest(query string, n int) []Suggestion {
	if n <= 0 {
		return nil
	}

	tokens := strings.Fields(strings.ToLower(query))
	if len(tokens) == 0 {
		if n >= len(builtins) {
			return All()
		}
		return All()[:n]
	}

	type scored struct {
		s     Suggestion
		score int
	}

	scores := make([]scored, len(builtins))
	for i, s := range builtins {
		label := strings.ToLower(s.Label)
		count := 0
		for _, tok := range tokens {
			if strings.Contains(label, tok) {
				count++
			}
		}
		scores[i] = scored{s: s, score: count}
	}

	// Stable selection sort for top-n by score (catalog order preserved for ties).
	results := make([]Suggestion, 0, n)
	used := make([]bool, len(scores))
	for range min(n, len(scores)) {
		best := -1
		for j, sc := range scores {
			if used[j] {
				continue
			}
			if sc.score == 0 {
				continue
			}
			if best == -1 || sc.score > scores[best].score {
				best = j
			}
		}
		if best == -1 {
			break
		}
		used[best] = true
		results = append(results, scores[best].s)
	}
	return results
}

// min is a small helper for Go versions before 1.21.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
