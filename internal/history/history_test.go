package history_test

import (
	"testing"
	"time"

	"github.com/your-org/cronparse/internal/history"
)

func ref(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestBefore_EveryMinute(t *testing.T) {
	base := ref("2024-06-01T12:05:00Z")
	times, err := history.Before("* * * * *", base, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 results, got %d", len(times))
	}
	expected := []string{
		"2024-06-01T12:04:00Z",
		"2024-06-01T12:03:00Z",
		"2024-06-01T12:02:00Z",
	}
	for i, ts := range times {
		if ts.UTC().Format(time.RFC3339) != expected[i] {
			t.Errorf("result[%d]: got %s, want %s", i, ts.UTC().Format(time.RFC3339), expected[i])
		}
	}
}

func TestBefore_SpecificMinute(t *testing.T) {
	// fires at minute 30 every hour
	base := ref("2024-06-01T14:00:00Z")
	times, err := history.Before("30 * * * *", base, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 2 {
		t.Fatalf("expected 2 results, got %d", len(times))
	}
	if times[0].UTC().Format(time.RFC3339) != "2024-06-01T13:30:00Z" {
		t.Errorf("first result: got %s", times[0].UTC().Format(time.RFC3339))
	}
	if times[1].UTC().Format(time.RFC3339) != "2024-06-01T12:30:00Z" {
		t.Errorf("second result: got %s", times[1].UTC().Format(time.RFC3339))
	}
}

func TestBefore_InvalidExpression(t *testing.T) {
	_, err := history.Before("invalid", time.Now(), 1)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestBefore_InvalidN(t *testing.T) {
	_, err := history.Before("* * * * *", time.Now(), 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestBefore_OrderedMostRecentFirst(t *testing.T) {
	base := ref("2024-06-01T10:00:00Z")
	times, err := history.Before("* * * * *", base, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(times); i++ {
		if !times[i-1].After(times[i]) {
			t.Errorf("results not ordered: times[%d]=%s is not after times[%d]=%s",
				i-1, times[i-1], i, times[i])
		}
	}
}
