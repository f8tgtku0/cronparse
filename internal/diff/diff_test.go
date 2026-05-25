package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/internal/diff"
)

func TestCompare_Identical(t *testing.T) {
	res, err := diff.Compare("0 9 * * 1", "0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Fields) != 0 {
		t.Errorf("expected no field diffs, got %d", len(res.Fields))
	}
	if !strings.Contains(res.Summary, "identical") {
		t.Errorf("expected 'identical' in summary, got: %s", res.Summary)
	}
}

func TestCompare_SingleFieldChange(t *testing.T) {
	res, err := diff.Compare("0 9 * * 1", "30 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Fields) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(res.Fields))
	}
	if res.Fields[0].Name != "minute" {
		t.Errorf("expected 'minute' diff, got %q", res.Fields[0].Name)
	}
	if res.Fields[0].From != "0" || res.Fields[0].To != "30" {
		t.Errorf("unexpected from/to: %q -> %q", res.Fields[0].From, res.Fields[0].To)
	}
}

func TestCompare_MultipleFieldChanges(t *testing.T) {
	res, err := diff.Compare("0 9 * * 1", "0 17 * 6 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Fields) != 3 {
		t.Errorf("expected 3 diffs, got %d", len(res.Fields))
	}
}

func TestCompare_InvalidFirst(t *testing.T) {
	_, err := diff.Compare("bad expr", "0 9 * * 1")
	if err == nil {
		t.Fatal("expected error for invalid first expression")
	}
	if !strings.Contains(err.Error(), "first expression") {
		t.Errorf("error should mention 'first expression': %v", err)
	}
}

func TestCompare_InvalidSecond(t *testing.T) {
	_, err := diff.Compare("0 9 * * 1", "not valid")
	if err == nil {
		t.Fatal("expected error for invalid second expression")
	}
	if !strings.Contains(err.Error(), "second expression") {
		t.Errorf("error should mention 'second expression': %v", err)
	}
}

func TestCompare_SummaryMentionsChangedFields(t *testing.T) {
	res, err := diff.Compare("*/5 * * * *", "*/10 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Summary, "minute") {
		t.Errorf("summary should mention 'minute': %s", res.Summary)
	}
}
