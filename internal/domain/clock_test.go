package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClockParseAndFormat(t *testing.T) {
	t.Parallel()

	clock, err := ParseClock("14:05:09")
	require.NoError(t, err)
	require.Equal(t, "14:05:09", clock.String())
}

func TestClockParseZero(t *testing.T) {
	t.Parallel()

	clock, err := ParseClock("00:00:00")
	require.NoError(t, err)
	require.True(t, clock.IsZero())
	require.Equal(t, "00:00:00", clock.String())
}

func TestClockParseInvalidFormats(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{"too few parts", "14:05"},
		{"too many parts", "14:05:09:01"},
		{"hours out of range", "25:00:00"},
		{"minutes out of range", "14:60:00"},
		{"seconds out of range", "14:05:60"},
		{"non-numeric hours", "ab:05:09"},
		{"non-numeric minutes", "14:ab:09"},
		{"non-numeric seconds", "14:05:ab"},
		{"single digit hours", "1:05:09"},
		{"single digit minutes", "14:5:09"},
		{"single digit seconds", "14:05:9"},
		{"negative number", "-1:05:09"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := ParseClock(tc.input)
			require.Error(t, err)
		})
	}
}

func TestClockComparisons(t *testing.T) {
	t.Parallel()

	early := mustParseClock(t, "10:00:00")
	late := mustParseClock(t, "12:00:00")
	same := mustParseClock(t, "10:00:00")

	require.True(t, early.Before(late))
	require.False(t, late.Before(early))
	require.False(t, early.Before(same))

	require.True(t, late.After(early))
	require.False(t, early.After(late))
	require.False(t, early.After(same))

	require.True(t, early.BeforeOrEqual(late))
	require.True(t, early.BeforeOrEqual(same))
	require.False(t, late.BeforeOrEqual(early))

	require.True(t, late.AfterOrEqual(early))
	require.True(t, early.AfterOrEqual(same))
	require.False(t, early.AfterOrEqual(late))
}

func TestClockAddHours(t *testing.T) {
	t.Parallel()

	clock := mustParseClock(t, "14:00:00")
	added := clock.AddHours(2)
	require.Equal(t, "16:00:00", added.String())
}

func TestDurationFormat(t *testing.T) {
	t.Parallel()

	start := mustParseClock(t, "14:40:00")
	end := mustParseClock(t, "15:04:00")

	duration := DurationBetween(start, end)
	require.Equal(t, "00:24:00", duration.String())
}

func TestDurationNegativeClampedToZero(t *testing.T) {
	t.Parallel()

	d := NewDuration(-10)
	require.Equal(t, 0, d.Seconds())
	require.True(t, d.IsZero())
}

func TestDurationAdd(t *testing.T) {
	t.Parallel()

	d1 := NewDuration(60)
	d2 := NewDuration(120)
	sum := d1.Add(d2)
	require.Equal(t, 180, sum.Seconds())
	require.Equal(t, "00:03:00", sum.String())
}

func TestDurationDiv(t *testing.T) {
	t.Parallel()

	d := NewDuration(600)
	require.Equal(t, 200, d.Div(3).Seconds())
	require.Equal(t, 0, d.Div(0).Seconds())
	require.Equal(t, 0, d.Div(-1).Seconds())
}

func TestDurationBetween(t *testing.T) {
	t.Parallel()

	start := mustParseClock(t, "14:00:00")
	end := mustParseClock(t, "14:00:00")

	d := DurationBetween(start, end)
	require.True(t, d.IsZero())
}

func mustParseClock(t *testing.T, value string) Clock {
	t.Helper()

	clock, err := ParseClock(value)
	require.NoError(t, err)

	return clock
}
