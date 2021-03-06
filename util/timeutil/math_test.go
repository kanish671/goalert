package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddClock(t *testing.T) {

	check := func(exp string, ts time.Time, c Clock) {
		t.Helper()
		res := AddClock(ts, c)
		assert.Equal(t, exp, res.String())
	}

	loc, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	check(
		"2018-03-11 01:00:00 -0600 CST",
		time.Date(2018, time.March, 11, 0, 0, 0, 0, loc),
		NewClock(1, 0),
	)
	check(
		"2018-03-11 03:00:00 -0500 CDT",
		time.Date(2018, time.March, 11, 0, 0, 0, 0, loc),
		NewClock(2, 0),
	)
	check(
		"2018-03-11 03:00:00 -0500 CDT",
		time.Date(2018, time.March, 11, 0, 0, 0, 0, loc),
		NewClock(2, 30),
	)

	check(
		"2017-01-03 12:00:00 -0600 CST",
		time.Date(2017, time.January, 1, 0, 0, 0, 0, loc),
		NewClock(60, 0),
	)

	check(
		"2020-11-01 01:30:00 -0500 CDT",
		time.Date(2020, 11, 2, 1, 30, 0, 0, loc),
		NewClock(-24, 0),
	)
}

func TestHoursBetween(t *testing.T) {
	check := func(exp int, a, b time.Time) {
		t.Helper()
		res := HoursBetween(a, b)
		assert.Equal(t, exp, res)
	}

	loc, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)
	// normal
	check(
		4,
		time.Date(2018, time.March, 9, 0, 0, 0, 0, loc),
		time.Date(2018, time.March, 9, 4, 0, 0, 0, loc),
	)

	// CST -> CDT
	check(
		4,
		time.Date(2018, time.March, 11, 0, 0, 0, 0, loc),
		time.Date(2018, time.March, 11, 4, 0, 0, 0, loc),
	)

	// CDT -> CST
	check(
		4,
		time.Date(2018, time.November, 4, 0, 0, 0, 0, loc),
		time.Date(2018, time.November, 4, 4, 0, 0, 0, loc),
	)

	check(
		49,
		time.Date(2020, 10, 30, 1, 30, 0, 0, loc),
		time.Date(2020, 11, 1, 2, 30, 0, 0, loc),
	)
}
