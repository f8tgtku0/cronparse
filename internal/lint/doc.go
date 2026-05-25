// Package lint analyses cron expressions and reports diagnostics.
//
// It distinguishes between hard errors (expressions that cannot be parsed)
// and warnings (expressions that are syntactically valid but may not behave
// as the author expects).
//
// Example usage:
//
//	diags := lint.Check("* * * * *")
//	for _, d := range diags {
//		fmt.Println(d)
//	}
//
// Current checks:
//   - Parse error — the expression is not a valid 5-field cron string.
//   - All-wildcards — every field is "*", which runs every minute.
//   - Dual day constraint — both day-of-month and day-of-week are non-wildcard;
//     most cron daemons combine them with OR semantics, which surprises users
//     who expect AND semantics.
package lint
