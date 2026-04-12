package analytics

import (
	"context"
)

type AnalyticsRepository interface {
	GetTotalCreatedLink(ctx context.Context, userID int64, filter FilterQuery) (int64, error)
	GetTotalClicks(ctx context.Context, userID int64, filter FilterQuery) (int64, error)
	GetUniqueClicks(ctx context.Context, userID int64, filter FilterQuery) (int64, error)
	GetClicks(ctx context.Context, linkID, userID int64, filter FilterQuery) (GetClicksResult, error)
	GetEngagementByDevices(ctx context.Context, linkID, userID int64, filter FilterQuery) ([]GetTotalEngagementByDeviceResult, error)
	GetEngagementByReferer(ctx context.Context, linkID, userID int64, filter FilterQuery) ([]GetTotalEngagementByRefererResult, error)
	GetEngagementByOvertime(ctx context.Context, linkID, userID int64, filter FilterQuery) ([]GetTotalEngagementByOvertimeResult, error)
}
