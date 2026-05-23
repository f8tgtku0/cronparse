package next

import (
	"time"

	"github.com/cronparse/internal/expr"
)

// After returns the next time the cron expression fires after t.
// It searches up to one year ahead; if no match is found it returns
// time.Time{} and false.
func After(e *expr.Expr, t time.Time) (time.Time, bool) {
	// Truncate to the next whole minute and advance by one minute so we
	// never return the same instant that was passed in.
	current := t.Truncate(time.Minute).Add(time.Minute)

	deadline := t.Add(366 * 24 * time.Hour)

	for !current.After(deadline) {
		if matches(e, current) {
			return current, true
		}
		current = current.Add(time.Minute)
	}

	return time.Time{}, false
}

// matches reports whether t satisfies every field of the expression.
func matches(e *expr.Expr, t time.Time) bool {
	if !e.Minute.Contains(t.Minute()) {
		return false
	}
	if !e.Hour.Contains(t.Hour()) {
		return false
	}
	if !e.DayOfMonth.Contains(t.Day()) {
		return false
	}
	if !e.Month.Contains(int(t.Month())) {
		return false
	}
	if !e.DayOfWeek.Contains(int(t.Weekday())) {
		return false
	}
	return true
}
