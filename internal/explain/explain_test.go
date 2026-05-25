package explain_test

import (
	"strings"
	"testing"

	"github.com/yourorg/cronparse/internal/explain"
)

func TestExplain_AllWildcards(t *testing.T) {
	res, err := explain.Explain("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expression != "* * * * *" {
		t.Errorf("expression mismatch: %q", res.Expression)
	}
	if len(res.Fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(res.Fields))
	}
	for _, f := range res.Fields {
		if !strings.HasPrefix(f.Meaning, "every ") {
			t.Errorf("field %q: expected 'every ...' got %q", f.Field, f.Meaning)
		}
	}
}

func TestExplain_SpecificValues(t *testing.T) {
	res, err := explain.Explain("30 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Fields[0].Meaning != "minute 30" {
		t.Errorf("minute: got %q", res.Fields[0].Meaning)
	}
	if res.Fields[1].Meaning != "hour 9" {
		t.Errorf("hour: got %q", res.Fields[1].Meaning)
	}
}

func TestExplain_StepExpression(t *testing.T) {
	res, err := explain.Explain("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Fields[0].Meaning, "every 5") {
		t.Errorf("step: got %q", res.Fields[0].Meaning)
	}
}

func TestExplain_RangeExpression(t *testing.T) {
	res, err := explain.Explain("0 9-17 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Fields[1].Meaning, "through") {
		t.Errorf("range: got %q", res.Fields[1].Meaning)
	}
}

func TestExplain_InvalidExpression(t *testing.T) {
	_, err := explain.Explain("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestExplain_SummaryContainsAllFields(t *testing.T) {
	res, err := explain.Explain("0 12 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Summary, "minute") || !strings.Contains(res.Summary, "hour") {
		t.Errorf("summary missing fields: %q", res.Summary)
	}
}
