package main

import (
	"strings"
	"testing"
)

func TestDiffCmd_IdenticalExpressions(t *testing.T) {
	out, code := integrationCmd(t, "diff", "0 9 * * 1", "0 9 * * 1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; output: %s", code, out)
	}
	if !strings.Contains(out, "identical") {
		t.Errorf("expected 'identical' in output, got: %s", out)
	}
}

func TestDiffCmd_ChangedMinute(t *testing.T) {
	out, code := integrationCmd(t, "diff", "0 9 * * 1", "30 9 * * 1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; output: %s", code, out)
	}
	if !strings.Contains(out, "minute") {
		t.Errorf("expected 'minute' in output, got: %s", out)
	}
	if !strings.Contains(out, "0") || !strings.Contains(out, "30") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}

func TestDiffCmd_InvalidExpression(t *testing.T) {
	_, code := integrationCmd(t, "diff", "bad expr here", "0 9 * * 1")
	if code == 0 {
		t.Fatal("expected non-zero exit for invalid expression")
	}
}

func TestDiffCmd_TooFewArgs(t *testing.T) {
	_, code := integrationCmd(t, "diff", "0 9 * * 1")
	if code != 2 {
		t.Fatalf("expected exit 2 for missing argument, got %d", code)
	}
}

func TestDiffCmd_NoArgs(t *testing.T) {
	_, code := integrationCmd(t, "diff")
	if code != 2 {
		t.Fatalf("expected exit 2 for no arguments, got %d", code)
	}
}
