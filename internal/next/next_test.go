package next_test

import (
	"testing"
	"time"

	"github.com/cronparse/internal/expr"
	"github.com/cronparse/internal/next"
)

func mustParse(t *testing.T, s string) *expr.Expr {
	t.Helper()
	e, err := expr.Parse(s)
	if err != nil {
		t.Fatalf("expr.Parse(%q): %v", s, err)
	}
	return e
}

func TestAfter_EveryMinute(t *testing.T) {
	e := mustParse(t, "* * * * *")
	base := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	got, ok := next.After(e, base)
	if !ok {
		t.Fatal("expected a next time, got none")
	}
	want := base.Add(time.Minute)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAfter_SpecificMinute(t *testing.T) {
	e := mustParse(t, "0 9 * * *")
	// base is 09:01 — next fire should be 09:00 the following day.
	base := time.Date(2024, 3, 10, 9, 1, 0, 0, time.UTC)
	got, ok := next.After(e, base)
	if !ok {
		t.Fatal("expected a next time, got none")
	}
	want := time.Date(2024, 3, 11, 9, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAfter_WeekdayFilter(t *testing.T) {
	// Every Monday at 08:00.
	e := mustParse(t, "0 8 * * 1")
	// 2024-01-15 is a Monday; base is 08:01 so next should be next Monday.
	base := time.Date(2024, 1, 15, 8, 1, 0, 0, time.UTC)
	got, ok := next.After(e, base)
	if !ok {
		t.Fatal("expected a next time, got none")
	}
	want := time.Date(2024, 1, 22, 8, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAfter_ExactlyOnMinuteBoundary(t *testing.T) {
	e := mustParse(t, "30 10 * * *")
	// base IS the fire time; After must return the NEXT occurrence.
	base := time.Date(2024, 5, 1, 10, 30, 0, 0, time.UTC)
	got, ok := next.After(e, base)
	if !ok {
		t.Fatal("expected a next time, got none")
	}
	want := time.Date(2024, 5, 2, 10, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
