package schedule_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/cronparse/internal/schedule"
)

const sampleJSON = `{
  "schedules": [
    {"name": "every-minute", "expression": "* * * * *"},
    {"name": "daily",        "expression": "0 9 * * *"}
  ]
}`

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "sched.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestLoadJSON_Valid(t *testing.T) {
	p := writeTempFile(t, sampleJSON)
	s := schedule.New()
	if err := s.LoadJSON(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}
}

func TestLoadJSON_MissingFile(t *testing.T) {
	s := schedule.New()
	if err := s.LoadJSON("/nonexistent/path.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadJSON_BadExpression(t *testing.T) {
	bad := `{"schedules":[{"name":"x","expression":"bad expr"}]}`
	p := writeTempFile(t, bad)
	s := schedule.New()
	if err := s.LoadJSON(p); err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestSaveJSON_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.json")

	s1 := schedule.New()
	_ = s1.Add("every-minute", "* * * * *")
	_ = s1.Add("weekly", "0 8 * * 1")

	if err := s1.SaveJSON(p); err != nil {
		t.Fatalf("SaveJSON: %v", err)
	}

	s2 := schedule.New()
	if err := s2.LoadJSON(p); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	if s2.Len() != s1.Len() {
		t.Fatalf("expected %d entries after round-trip, got %d", s1.Len(), s2.Len())
	}
}
