package domain

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	secondsInMinute = 60
	secondsInHour   = 60 * secondsInMinute
	secondsInDay    = 24 * secondsInHour
)

type Clock struct {
	seconds int
}

type Duration struct {
	seconds int
}

func ParseClock(value string) (Clock, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return Clock{}, fmt.Errorf("invalid time %q", value)
	}

	hours, err := parseTimePart(parts[0], 23)
	if err != nil {
		return Clock{}, fmt.Errorf("invalid hours in %q: %w", value, err)
	}

	minutes, err := parseTimePart(parts[1], 59)
	if err != nil {
		return Clock{}, fmt.Errorf("invalid minutes in %q: %w", value, err)
	}

	seconds, err := parseTimePart(parts[2], 59)
	if err != nil {
		return Clock{}, fmt.Errorf("invalid seconds in %q: %w", value, err)
	}

	return Clock{seconds: hours*secondsInHour + minutes*secondsInMinute + seconds}, nil
}

func NewDuration(seconds int) Duration {
	if seconds < 0 {
		seconds = 0
	}

	return Duration{seconds: seconds}
}

func DurationBetween(start, end Clock) Duration {
	return NewDuration(end.seconds - start.seconds)
}

func (c Clock) AddHours(hours int) Clock {
	return Clock{seconds: c.seconds + hours*secondsInHour}
}

func (c Clock) Before(other Clock) bool {
	return c.seconds < other.seconds
}

func (c Clock) After(other Clock) bool {
	return c.seconds > other.seconds
}

func (c Clock) BeforeOrEqual(other Clock) bool {
	return c.seconds <= other.seconds
}

func (c Clock) AfterOrEqual(other Clock) bool {
	return c.seconds >= other.seconds
}

func (c Clock) IsZero() bool {
	return c.seconds == 0
}

func (c Clock) String() string {
	seconds := c.seconds % secondsInDay
	if seconds < 0 {
		seconds += secondsInDay
	}

	return formatHMS(seconds)
}

func (d Duration) Seconds() int {
	return d.seconds
}

func (d Duration) IsZero() bool {
	return d.seconds == 0
}

func (d Duration) Add(other Duration) Duration {
	return Duration{seconds: d.seconds + other.seconds}
}

func (d Duration) Div(divisor int) Duration {
	if divisor <= 0 {
		return Duration{}
	}

	return Duration{seconds: d.seconds / divisor}
}

func (d Duration) String() string {
	return formatHMS(d.seconds)
}

func formatHMS(totalSeconds int) string {
	hours := totalSeconds / secondsInHour
	minutes := (totalSeconds % secondsInHour) / secondsInMinute
	seconds := totalSeconds % secondsInMinute

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func parseTimePart(value string, max int) (int, error) {
	if len(value) != 2 {
		return 0, fmt.Errorf("must contain exactly two digits")
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if number < 0 || number > max {
		return 0, fmt.Errorf("must be between 0 and %d", max)
	}

	return number, nil
}
