// Package template provides cron expression templates for common scheduling
// patterns, allowing users to generate expressions from named presets.
package template

import (
	"fmt"
	"sort"
	"strings"
)

// Template represents a named cron expression preset with a description.
type Template struct {
	// Name is the short identifier for the template (e.g. "hourly").
	Name string
	// Expression is the standard 5-field cron expression.
	Expression string
	// Description is a human-readable summary of when the expression fires.
	Description string
}

// builtins holds the predefined cron templates keyed by name.
var builtins = map[string]Template{
	"yearly": {
		Name:        "yearly",
		Expression:  "0 0 1 1 *",
		Description: "Once a year at midnight on January 1st",
	},
	"annually": {
		Name:        "annually",
		Expression:  "0 0 1 1 *",
		Description: "Once a year at midnight on January 1st (alias for yearly)",
	},
	"monthly": {
		Name:        "monthly",
		Expression:  "0 0 1 * *",
		Description: "Once a month at midnight on the 1st",
	},
	"weekly": {
		Name:        "weekly",
		Expression:  "0 0 * * 0",
		Description: "Once a week at midnight on Sunday",
	},
	"daily": {
		Name:        "daily",
		Expression:  "0 0 * * *",
		Description: "Once a day at midnight",
	},
	"midnight": {
		Name:        "midnight",
		Expression:  "0 0 * * *",
		Description: "Once a day at midnight (alias for daily)",
	},
	"hourly": {
		Name:        "hourly",
		Expression:  "0 * * * *",
		Description: "Once an hour at the start of the hour",
	},
	"every-minute": {
		Name:        "every-minute",
		Expression:  "* * * * *",
		Description: "Every minute",
	},
	"every-5-minutes": {
		Name:        "every-5-minutes",
		Expression:  "*/5 * * * *",
		Description: "Every 5 minutes",
	},
	"every-15-minutes": {
		Name:        "every-15-minutes",
		Expression:  "*/15 * * * *",
		Description: "Every 15 minutes",
	},
	"every-30-minutes": {
		Name:        "every-30-minutes",
		Expression:  "*/30 * * * *",
		Description: "Every 30 minutes",
	},
	"weekdays": {
		Name:        "weekdays",
		Expression:  "0 9 * * 1-5",
		Description: "At 9:00 AM every weekday (Monday through Friday)",
	},
	"weekends": {
		Name:        "weekends",
		Expression:  "0 9 * * 6,0",
		Description: "At 9:00 AM on Saturday and Sunday",
	},
}

// Lookup returns the Template for the given name, or an error if not found.
// Name matching is case-insensitive.
func Lookup(name string) (Template, error) {
	t, ok := builtins[strings.ToLower(name)]
	if !ok {
		return Template{}, fmt.Errorf("unknown template %q; use List() to see available templates", name)
	}
	return t, nil
}

// List returns all built-in templates sorted alphabetically by name.
func List() []Template {
	result := make([]Template, 0, len(builtins))
	for _, t := range builtins {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// Names returns a sorted slice of all available template names.
func Names() []string {
	names := make([]string, 0, len(builtins))
	for name := range builtins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
