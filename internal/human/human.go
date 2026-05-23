// Package human provides human-readable descriptions of cron expressions.
package human

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronparse/internal/expr"
)

// Describe returns a human-readable description of a parsed cron expression.
func Describe(e *expr.Expr) string {
	parts := []string{}

	minDesc := describeField(e.Minute, "minute", "every minute", "minute %s")
	hourDesc := describeField(e.Hour, "hour", "every hour", "hour %s")
	dayDesc := describeField(e.DayOfMonth, "day of month", "every day", "day %s of the month")
	monthDesc := describeMonth(e.Month)
	weekDesc := describeWeekday(e.DayOfWeek)

	if isWildcard(e.Minute) && isWildcard(e.Hour) && isWildcard(e.DayOfMonth) && isWildcard(e.Month) && isWildcard(e.DayOfWeek) {
		return "every minute"
	}

	if !isWildcard(e.Minute) {
		parts = append(parts, minDesc)
	}
	if !isWildcard(e.Hour) {
		parts = append(parts, hourDesc)
	}
	if !isWildcard(e.DayOfMonth) {
		parts = append(parts, dayDesc)
	}
	if !isWildcard(e.Month) {
		parts = append(parts, monthDesc)
	}
	if !isWildcard(e.DayOfWeek) {
		parts = append(parts, weekDesc)
	}

	if len(parts) == 0 {
		return "every minute"
	}
	return "At " + strings.Join(parts, ", ")
}

func isWildcard(raw string) bool {
	return raw == "*"
}

func describeField(raw, _ string, wildcardDesc, fmtStr string) string {
	if isWildcard(raw) {
		return wildcardDesc
	}
	return fmt.Sprintf(fmtStr, raw)
}

var monthNames = map[string]string{
	"1": "January", "2": "February", "3": "March", "4": "April",
	"5": "May", "6": "June", "7": "July", "8": "August",
	"9": "September", "10": "October", "11": "November", "12": "December",
}

var weekdayNames = map[string]string{
	"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
	"4": "Thursday", "5": "Friday", "6": "Saturday",
}

func describeMonth(raw string) string {
	if isWildcard(raw) {
		return "every month"
	}
	if name, ok := monthNames[raw]; ok {
		return "in " + name
	}
	return "in month " + raw
}

func describeWeekday(raw string) string {
	if isWildcard(raw) {
		return "every weekday"
	}
	if name, ok := weekdayNames[raw]; ok {
		return "on " + name
	}
	return "on weekday " + raw
}
