// Package schedule provides multi-expression scheduling support for cronparse.
//
// A Schedule holds a named collection of cron expressions. Each expression is
// parsed and validated on insertion via [Schedule.Add]. Once populated, the
// schedule can be queried for the next run time of every entry relative to an
// arbitrary reference time using [Schedule.NextAll], which returns results
// sorted in ascending chronological order.
//
// Schedules can be persisted to and loaded from JSON files with
// [Schedule.SaveJSON] and [Schedule.LoadJSON]. The expected file format is:
//
//	{
//	  "schedules": [
//	    { "name": "my-job", "expression": "0 9 * * 1-5" }
//	  ]
//	}
//
// Example usage:
//
//	s := schedule.New()
//	_ = s.Add("daily-report", "0 8 * * *")
//	_ = s.Add("hourly-sync",  "0 * * * *")
//	for _, ne := range s.NextAll(time.Now()) {
//	    fmt.Printf("%s → %s\n", ne.Name, ne.Next.Format(time.RFC3339))
//	}
package schedule
