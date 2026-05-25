package schedule_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronparse/internal/schedule"
)

func TestAdd_Valid(t *testing.T) {
	s := schedule.New()
	if err := s.Add("every-minute", "* * * * *"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", s.Len())
	}
}

func TestAdd_Invalid(t *testing.T) {
	s := schedule.New()
	if err := s.Add("bad", "not a cron"); err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestNextAll_Ordering(t *testing.T) {
	s := schedule.New()
	// hourly at :30 vs every minute — every-minute fires sooner
	_ = s.Add("hourly", "30 * * * *")
	_ = s.Add("every-minute", "* * * * *")

	// Use a time that is NOT on the :30 mark so ordering is deterministic.
	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	results := s.NextAll(base)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Next.Before(results[1].Next) && results[0].Next != results[1].Next {
		t.Errorf("results not sorted: %v >= %v", results[0].Next, results[1].Next)
	}
	if results[0].Name != "every-minute" {
		t.Errorf("expected every-minute first, got %q", results[0].Name)
	}
}

func TestNextAll_Empty(t *testing.T) {
	s := schedule.New()
	results := s.NextAll(time.Now())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestNextAll_MultipleEntries(t *testing.T) {
	s := schedule.New()
	_ = s.Add("midnight", "0 0 * * *")
	_ = s.Add("noon", "0 12 * * *")
	_ = s.Add("every-minute", "* * * * *")

	base := time.Date(2024, 6, 1, 6, 0, 0, 0, time.UTC)
	results := s.NextAll(base)

	for i := 1; i < len(results); i++ {
		if results[i].Next.Before(results[i-1].Next) {
			t.Errorf("entry %d (%s) fires before entry %d (%s)",
				i, results[i].Name, i-1, results[i-1].Name)
		}
	}
}
