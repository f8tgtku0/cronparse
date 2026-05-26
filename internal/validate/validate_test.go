package validate_test

import (
	"testing"

	"github.com/yourorg/cronparse/internal/validate"
)

func TestValidate_AllWildcards(t *testing.T) {
	result := validate.Validate("* * * * *")
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidate_SpecificSchedule(t *testing.T) {
	result := validate.Validate("30 9 * * 1-5")
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidate_InvalidMinute(t *testing.T) {
	result := validate.Validate("60 * * * *")
	if result.Valid {
		t.Error("expected invalid for minute=60")
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error")
	}
}

func TestValidate_InvalidHour(t *testing.T) {
	result := validate.Validate("0 25 * * *")
	if result.Valid {
		t.Error("expected invalid for hour=25")
	}
}

func TestValidate_InvalidDayOfMonth(t *testing.T) {
	result := validate.Validate("0 0 32 * *")
	if result.Valid {
		t.Error("expected invalid for day=32")
	}
}

func TestValidate_InvalidMonth(t *testing.T) {
	result := validate.Validate("0 0 1 13 *")
	if result.Valid {
		t.Error("expected invalid for month=13")
	}
}

func TestValidate_InvalidWeekday(t *testing.T) {
	result := validate.Validate("0 0 * * 8")
	if result.Valid {
		t.Error("expected invalid for weekday=8")
	}
}

func TestValidate_TooFewFields(t *testing.T) {
	result := validate.Validate("* * * *")
	if result.Valid {
		t.Error("expected invalid for too few fields")
	}
}

func TestValidate_TooManyFields(t *testing.T) {
	result := validate.Validate("* * * * * *")
	if result.Valid {
		t.Error("expected invalid for too many fields")
	}
}

func TestValidate_StepExpression(t *testing.T) {
	result := validate.Validate("*/5 * * * *")
	if !result.Valid {
		t.Errorf("expected valid step expression, got errors: %v", result.Errors)
	}
}

func TestValidate_RangeExpression(t *testing.T) {
	result := validate.Validate("0 9-17 * * 1-5")
	if !result.Valid {
		t.Errorf("expected valid range expression, got errors: %v", result.Errors)
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	// range where start > end
	result := validate.Validate("0 17-9 * * *")
	if result.Valid {
		t.Error("expected invalid for reversed range 17-9")
	}
}

func TestValidate_BothDayFields_Warning(t *testing.T) {
	result := validate.Validate("0 0 15 * 1")
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
	if len(result.Warnings) == 0 {
		t.Error("expected a warning when both day-of-month and day-of-week are set")
	}
}

func TestValidate_ZeroStepError(t *testing.T) {
	result := validate.Validate("*/0 * * * *")
	if result.Valid {
		t.Error("expected invalid for step of zero")
	}
}

func TestResult_Summary(t *testing.T) {
	result := validate.Validate("60 * * * *")
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary for invalid expression")
	}
}

func TestResult_Summary_Valid(t *testing.T) {
	result := validate.Validate("* * * * *")
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary even for valid expression")
	}
}
