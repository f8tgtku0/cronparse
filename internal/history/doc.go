// Package history provides functionality for computing past execution times
// of a cron expression relative to a given reference time.
//
// Given a parsed cron expression and a reference time, Before returns the N
// most recent times the expression would have fired before that instant.
// Results are returned in reverse-chronological order (most recent first).
//
// Example usage:
//
//	expr, _ := expr.Parse("0 9 * * 1-5")
//	times, err := history.Before(expr, time.Now(), 5)
//	for _, t := range times {
//		fmt.Println(t)
//	}
package history
