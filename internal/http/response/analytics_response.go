package response

import "time"

type GetClicksResponse struct {
	TotalClick  int64 `json:"totalClick"`
	UniqueClick int64 `json:"uniqueClick"`
}

type GetTotalEngagementByDeviceResponse struct {
	Device string `json:"device"`
	Clicks int64  `json:"clicks"`
}

type GetTotalEngagementByRefererResponse struct {
	Referer string `json:"referer"`
	Clicks  int64  `json:"clicks"`
}

type GetTotalEngagementByOvertimeResponse struct {
	Date   time.Time `json:"date"`
	Clicks int64     `json:"clicks"`
}
