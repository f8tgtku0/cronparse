package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_AllWildcards(t *testing.T) {
	e, err := Parse("* * * * *")
	require.NoError(t, err)
	assert.Equal(t, 60, len(e.Minute.Values))
	assert.Equal(t, 24, len(e.Hour.Values))
	assert.Equal(t, 31, len(e.DayOfMonth.Values))
	assert.Equal(t, 12, len(e.Month.Values))
	assert.Equal(t, 7, len(e.DayOfWeek.Values))
}

func TestParse_SpecificSchedule(t *testing.T) {
	// Every day at 02:30
	e, err := Parse("30 2 * * *")
	require.NoError(t, err)
	assert.Equal(t, []int{30}, e.Minute.Values)
	assert.Equal(t, []int{2}, e.Hour.Values)
}

func TestParse_WeekdaySchedule(t *testing.T) {
	// Mon-Fri at midnight
	e, err := Parse("0 0 * * 1-5")
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, e.DayOfWeek.Values)
}

func TestParse_TooFewFields(t *testing.T) {
	_, err := Parse("* * * *")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected 5 fields")
}

func TestParse_TooManyFields(t *testing.T) {
	_, err := Parse("* * * * * *")
	assert.Error(t, err)
}

func TestParse_InvalidFieldPropagates(t *testing.T) {
	_, err := Parse("60 * * * *")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse error in field 1")
}

func TestParse_Raw(t *testing.T) {
	raw := "*/5 * * * *"
	e, err := Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, raw, e.Raw)
}

func TestExpression_String(t *testing.T) {
	e, err := Parse("0 12 * * 0")
	require.NoError(t, err)
	s := e.String()
	assert.Contains(t, s, "0 12 * * 0")
	assert.Contains(t, s, "minute=")
}
