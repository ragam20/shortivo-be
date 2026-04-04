package utils

import "time"

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
