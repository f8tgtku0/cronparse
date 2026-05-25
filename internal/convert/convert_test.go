package convert_test

import (
	"testing"

	"github.com/yourorg/cronparse/internal/convert"
)

func TestConvert_SameDialect(t *testing.T) {
	res, err := convert.Convert("0 9 * * 1", convert.DialectStandard, convert.DialectStandard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expression != "0 9 * * 1" {
		t.Errorf("expected unchanged expression, got %q", res.Expression)
	}
}

func TestConvert_StandardToQuartz(t *testing.T) {
	res, err := convert.Convert("30 8 * * *", convert.DialectStandard, convert.DialectQuartz)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "0 30 8 * * * *"
	if res.Expression != want {
		t.Errorf("got %q, want %q", res.Expression, want)
	}
	if len(res.Notes) == 0 {
		t.Error("expected at least one note about added fields")
	}
}

func TestConvert_QuartzToStandard_ZeroSeconds(t *testing.T) {
	res, err := convert.Convert("0 30 8 * * *", convert.DialectQuartz, convert.DialectStandard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "30 8 * * *"
	if res.Expression != want {
		t.Errorf("got %q, want %q", res.Expression, want)
	}
}

func TestConvert_QuartzToStandard_NonZeroSeconds_Note(t *testing.T) {
	res, err := convert.Convert("30 0 9 * * *", convert.DialectQuartz, convert.DialectStandard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Notes) == 0 {
		t.Error("expected a note about dropped seconds field")
	}
}

func TestConvert_StandardToEventBridge(t *testing.T) {
	res, err := convert.Convert("0 12 1 1 *", convert.DialectStandard, convert.DialectEventBridge)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "0 12 1 1 * *"
	if res.Expression != want {
		t.Errorf("got %q, want %q", res.Expression, want)
	}
}

func TestConvert_EventBridgeToStandard_YearNote(t *testing.T) {
	res, err := convert.Convert("0 12 1 1 * 2025", convert.DialectEventBridge, convert.DialectStandard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "0 12 1 1 *"
	if res.Expression != want {
		t.Errorf("got %q, want %q", res.Expression, want)
	}
	if len(res.Notes) == 0 {
		t.Error("expected a note about dropped year field")
	}
}

func TestConvert_InvalidFieldCount(t *testing.T) {
	_, err := convert.Convert("* *", convert.DialectStandard, convert.DialectQuartz)
	if err == nil {
		t.Error("expected error for too few fields")
	}
}

func TestConvert_UnknownDialect(t *testing.T) {
	_, err := convert.Convert("* * * * *", convert.DialectStandard, "unknown")
	if err == nil {
		t.Error("expected error for unknown destination dialect")
	}
}
