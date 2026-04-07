package utils

import (
	"fmt"
	"time"
)

func ParseDateLocation(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, fmt.Errorf("invalid time format, must be RFC3339")
	}

	return &t, nil
}

func ToTimeRFC3339(tm time.Time) string {
	return tm.Format(time.RFC3339)
}

type DateRange struct {
	From time.Time
	To   time.Time
}

func GetDateRange(t time.Time) DateRange {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return DateRange{
		From: start,
		To:   end,
	}
}
