package lint_test

import (
	"testing"

	"github.com/yourorg/cronparse/internal/lint"
)

func TestCheck_ParseError(t *testing.T) {
	diags := lint.Check("not a cron")
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic for invalid expression")
	}
	if diags[0].Severity != lint.Error {
		t.Errorf("expected Error severity, got %s", diags[0].Severity)
	}
}

func TestCheck_AllWildcards(t *testing.T) {
	diags := lint.Check("* * * * *")
	if len(diags) == 0 {
		t.Fatal("expected warning for all-wildcard expression")
	}
	for _, d := range diags {
		if d.Severity != lint.Warning {
			t.Errorf("expected Warning, got %s: %s", d.Severity, d.Message)
		}
	}
}

func TestCheck_BothDayFields(t *testing.T) {
	diags := lint.Check("0 12 15 * 1")
	if len(diags) == 0 {
		t.Fatal("expected warning when both day fields are set")
	}
	found := false
	for _, d := range diags {
		if d.Severity == lint.Warning {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one Warning diagnostic")
	}
}

func TestCheck_CleanExpression(t *testing.T) {
	diags := lint.Check("30 9 * * 1")
	if len(diags) != 0 {
		t.Errorf("expected no diagnostics for clean expression, got %v", diags)
	}
}

func TestDiagnosticString(t *testing.T) {
	d := lint.Diagnostic{Severity: lint.Warning, Message: "test message"}
	if d.String() != "warning: test message" {
		t.Errorf("unexpected string: %s", d.String())
	}
}
