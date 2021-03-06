package schedule

import (
	"sort"
	"strings"
	"time"
)

// A FixedShift represents an on-call user with a start and end time.
type FixedShift struct {
	Start, End time.Time
	UserID     string
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func clampShiftTimes(start, end time.Time, shifts []FixedShift) []FixedShift {
	result := shifts[:0]
	// trim/clamp shift times
	for _, s := range shifts {
		if s.Start.Before(start) {
			s.Start = start
		}
		if s.End.After(end) {
			s.End = end
		}
		if !s.End.After(s.Start) {
			continue
		}

		result = append(result, s)
	}
	return result
}

func mergeShiftsByTime(shifts []FixedShift) []FixedShift {
	if len(shifts) == 0 {
		return shifts
	}

	sort.Slice(shifts, func(i, j int) bool { return shifts[i].Start.Before(shifts[j].Start) })
	result := shifts[:1]
	for _, s := range shifts[1:] {
		l := len(result) - 1

		if !s.End.After(s.Start) {
			// omit empty time range
			continue
		}
		if s.Start.After(result[l].End) {
			result = append(result, s)
			continue
		}

		// TODO: remove once we switch to uuid.UUID
		s.UserID = strings.ToLower(s.UserID)

		result[l].End = maxTime(result[l].End, s.End)
	}

	return result
}
func mergeShifts(shifts []FixedShift) []FixedShift {
	m := make(map[string][]FixedShift)
	for _, s := range shifts {
		m[s.UserID] = append(m[s.UserID], s)
	}
	result := shifts[:0]
	for _, s := range m {
		result = append(result, mergeShiftsByTime(s)...)
	}

	// Return deterministic output by sorting by Start,
	// or UserID if the Start times are equal.
	sort.Slice(result, func(i, j int) bool {
		if !result[i].Start.Equal(result[j].Start) {
			return result[i].Start.Before(result[j].Start)
		}
		return result[i].UserID < result[j].UserID
	})

	return result
}
