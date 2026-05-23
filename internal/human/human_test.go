package human_test

import (
	"testing"

	"github.com/yourorg/cronparse/internal/expr"
	"github.com/yourorg/cronparse/internal/human"
)

func mustParse(t *testing.T, s string) *expr.Expr {
	t.Helper()
	e, err := expr.Parse(s)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", s, err)
	}
	return e
}

func TestDescribe_AllWildcards(t *testing.T) {
	e := mustParse(t, "* * * * *")
	got := human.Describe(e)
	if got != "every minute" {
		t.Errorf("expected 'every minute', got %q", got)
	}
}

func TestDescribe_SpecificMinute(t *testing.T) {
	e := mustParse(t, "30 * * * *")
	got := human.Describe(e)
	want := "At minute 30"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDescribe_SpecificHourAndMinute(t *testing.T) {
	e := mustParse(t, "0 9 * * *")
	got := human.Describe(e)
	want := "At minute 0, hour 9"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDescribe_MonthName(t *testing.T) {
	e := mustParse(t, "0 0 1 12 *")
	got := human.Describe(e)
	want := "At minute 0, hour 0, day 1 of the month, in December"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDescribe_WeekdayName(t *testing.T) {
	e := mustParse(t, "0 9 * * 1")
	got := human.Describe(e)
	want := "At minute 0, hour 9, on Monday"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestDescribe_StepExpression(t *testing.T) {
	e := mustParse(t, "*/15 * * * *")
	got := human.Describe(e)
	want := "At minute */15"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
