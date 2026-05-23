// Package main provides the cronparse command-line tool.
//
// cronparse parses and validates human-readable cron expressions,
// describes them in plain English, and predicts the next N scheduled
// run times from a given reference point.
//
// Usage:
//
//	cronparse [flags] <cron expression>
//
// Flags:
//
//	-n int
//		number of next run times to display (default 5)
//	-from string
//		start time in RFC3339 format (default: now)
//
// Examples:
//
//	# Show next 5 runs of a weekday morning job
//	cronparse "0 9 * * 1-5"
//
//	# Show next 10 runs of an every-5-minute job starting from a fixed time
//	cronparse -n 10 -from 2024-06-01T00:00:00Z "*/5 * * * *"
//
//	# Monthly job on the 1st at midnight
//	cronparse "0 0 1 * *"
package main
