package main

import (
	"os/exec"
	"strings"
	"testing"
)

// integrationCmd builds the binary and runs it with the given args,
// returning combined stdout+stderr output.
func integrationCmd(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		}
	}
	return string(out), code
}

func TestMain_NoArgs(t *testing.T) {
	out, code := integrationCmd(t)
	if code == 0 {
		t.Fatal("expected non-zero exit code when no args given")
	}
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage message, got: %s", out)
	}
}

func TestMain_InvalidExpression(t *testing.T) {
	out, code := integrationCmd(t, "not a valid cron")
	if code == 0 {
		t.Fatal("expected non-zero exit code for invalid expression")
	}
	if !strings.Contains(out, "error:") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestMain_ValidExpression(t *testing.T) {
	out, code := integrationCmd(t,
		"-n", "3",
		"-from", "2024-01-15T09:00:00Z",
		"0 9 * * 1-5",
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; output: %s", code, out)
	}
	if !strings.Contains(out, "Expression") {
		t.Errorf("expected Expression line, got: %s", out)
	}
	if !strings.Contains(out, "Description") {
		t.Errorf("expected Description line, got: %s", out)
	}
	if !strings.Contains(out, "Next 3 run(s)") {
		t.Errorf("expected next runs header, got: %s", out)
	}
}

func TestMain_InvalidFromFlag(t *testing.T) {
	out, code := integrationCmd(t, "-from", "not-a-date", "* * * * *")
	if code == 0 {
		t.Fatal("expected non-zero exit code for invalid --from")
	}
	if !strings.Contains(out, "error:") {
		t.Errorf("expected error message, got: %s", out)
	}
}
