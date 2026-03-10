package repo

import (
	"context"
	"time"

	"shortivo/internal/app/analytics"
	"shortivo/internal/database/entity"

	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) analytics.AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetClicks(ctx context.Context, linkID, userID int64) (analytics.GetClicksResult, error) {
	var link entity.LinkEntity
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", linkID, userID).
		First(&link).Error; err != nil {
		return analytics.GetClicksResult{}, err
	}

	return analytics.GetClicksResult{
		TotalClick:  link.ClickCount,
		UniqueClick: link.UniqueClick,
	}, nil
}

func (r *analyticsRepository) GetEngagementByDevices(ctx context.Context, linkID int64) ([]analytics.GetTotalEngagementByDeviceResult, error) {
	var result []analytics.GetTotalEngagementByDeviceResult

	if err := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("Device, COUNT(*) AS clicks").
		Where("link_id = ?", linkID).
		Group("device").
		Order("clicks DESC").
		Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *analyticsRepository) GetEngagementByReferer(ctx context.Context, linkID int64) ([]analytics.GetTotalEngagementByRefererResult, error) {
	var result []analytics.GetTotalEngagementByRefererResult

	if err := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("COALESCE(NULLIF(referer, ''), 'direct') AS referer, COUNT(*) AS clicks").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(referer, ''), 'direct')").
		Order("clicks DESC").
		Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *analyticsRepository) GetEngagementByOvertime(ctx context.Context, linkID int64, startTime, endTime time.Time) ([]analytics.GetTotalEngagementByOvertimeResult, error) {
	var result []analytics.GetTotalEngagementByOvertimeResult

	if err := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("Date(created_at) AS date, COUNT(*) AS clicks").
		Where("link_id = ? AND created_at >= ? AND created_at <= ?", linkID, startTime, endTime).
		Group("created_at").
		Order("date DESC").
		Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
