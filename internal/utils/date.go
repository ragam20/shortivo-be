package utils

import "time"

func ToTimeRFC3339(tm time.Time) string {
	return tm.Format(time.RFC3339)
}
