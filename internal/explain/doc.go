// Package explain provides structured, field-by-field explanations of cron
// expressions.
//
// Each field (minute, hour, day-of-month, month, day-of-week) is annotated
// with its raw token and a human-readable meaning. The result also includes
// a one-line summary joining all field descriptions.
//
// Example:
//
//	res, err := explain.Explain("*/15 9-17 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, f := range res.Fields {
//		fmt.Printf("%-15s %-10s %s\n", f.Field, f.Raw, f.Meaning)
//	}
package explain
