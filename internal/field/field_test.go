package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Wildcard(t *testing.T) {
	f, err := Parse("*", Minute)
	require.NoError(t, err)
	assert.Equal(t, 60, len(f.Values))
	assert.Equal(t, 0, f.Values[0])
	assert.Equal(t, 59, f.Values[59])
}

func TestParse_Literal(t *testing.T) {
	f, err := Parse("5", Hour)
	require.NoError(t, err)
	assert.Equal(t, []int{5}, f.Values)
}

func TestParse_LiteralOutOfBounds(t *testing.T) {
	_, err := Parse("60", Minute)
	assert.Error(t, err)
}

func TestParse_Range(t *testing.T) {
	f, err := Parse("1-5", DayOfWeek)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, f.Values)
}

func TestParse_Step(t *testing.T) {
	f, err := Parse("*/15", Minute)
	require.NoError(t, err)
	assert.Equal(t, []int{0, 15, 30, 45}, f.Values)
}

func TestParse_StepWithRange(t *testing.T) {
	f, err := Parse("1-6/2", Month)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, f.Values)
}

func TestParse_List(t *testing.T) {
	f, err := Parse("1,15,30", Minute)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 15, 30}, f.Values)
}

func TestParse_ListWithRange(t *testing.T) {
	f, err := Parse("1-3,5", DayOfMonth)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 5}, f.Values)
}

func TestParse_InvalidStep(t *testing.T) {
	_, err := Parse("*/0", Minute)
	assert.Error(t, err)
}

func TestParse_InvalidRange(t *testing.T) {
	_, err := Parse("5-3", Hour)
	assert.Error(t, err)
}

func TestParse_InvalidLiteral(t *testing.T) {
	_, err := Parse("abc", Minute)
	assert.Error(t, err)
}
