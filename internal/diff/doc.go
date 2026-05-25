// Package diff provides comparison utilities for cron expressions.
//
// Given two valid cron expression strings, Compare identifies which fields
// differ (minute, hour, day-of-month, month, day-of-week) and returns a
// human-readable summary of the changes.
//
// Example:
//
//	res, err := diff.Compare("0 9 * * 1", "30 17 * * 5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Summary)
//	for _, f := range res.Fields {
//		fmt.Printf("  %s: %s -> %s\n", f.Name, f.From, f.To)
//	}
package diff
