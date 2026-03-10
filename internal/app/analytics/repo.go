package analytics

import (
	"context"
	"time"
)

type AnalyticsRepository interface {
	GetClicks(ctx context.Context, linkID, userID int64) (GetClicksResult, error)
	GetEngagementByDevices(ctx context.Context, linkID int64) ([]GetTotalEngagementByDeviceResult, error)
	GetEngagementByReferer(ctx context.Context, linkID int64) ([]GetTotalEngagementByRefererResult, error)
	GetEngagementByOvertime(ctx context.Context, linkID int64, startTime, endTime time.Time) ([]GetTotalEngagementByOvertimeResult, error)
}
