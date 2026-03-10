package analytics

import "time"

type GetTotalClicksQuery struct {
	UserID int64
	LinkID int64
}

type GetTotalEngagementByOvertimeQuery struct {
	UserID    int64
	LinkID    int64
	StartTime time.Time
	EndTime   time.Time
}

type GetClicksResult struct {
	TotalClick  int64
	UniqueClick int64
}

type GetTotalEngagementByDeviceResult struct {
	Device string
	Clicks int64
}

type GetTotalEngagementByRefererResult struct {
	Referer string
	Clicks  int64
}

type GetTotalEngagementByOvertimeResult struct {
	Date   time.Time
	Clicks int64
}
