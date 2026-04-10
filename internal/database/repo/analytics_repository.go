package repo

import (
	"context"

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

func (r *analyticsRepository) GetTotalClicks(ctx context.Context, userID int64, filter analytics.FilterQuery) (int64, error) {
	var totalClicks int64

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Joins("JOIN links l ON l.id = link_clicks.link_id").
		Where("l.user_id = ? AND l.deleted_at IS NULL", userID)

	db = applyAnalyticsDateFilter(db, "link_clicks.created_at", filter)

	if err := db.Count(&totalClicks).Error; err != nil {
		return 0, err
	}

	return totalClicks, nil
}

func (r *analyticsRepository) GetTotalCreatedLink(ctx context.Context, userID int64, filter analytics.FilterQuery) (int64, error) {
	var result int64

	db := r.db.WithContext(ctx).
		Model(&entity.LinkEntity{}).
		Where("user_id = ?", userID)

	db = applyAnalyticsDateFilter(db, "created_at", filter)

	if err := db.Count(&result).Error; err != nil {
		return 0, err
	}

	return result, nil
}

func (r *analyticsRepository) GetUniqueClicks(ctx context.Context, userID int64, filter analytics.FilterQuery) (int64, error) {
	var uniqueClicks int64

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL", userID)

	db = applyAnalyticsDateFilter(db, "link_clicks.created_at", filter)

	if err := db.
		Select("COUNT(DISTINCT (link_clicks.ip_address, link_clicks.user_agent))").
		Scan(&uniqueClicks).Error; err != nil {
		return 0, err
	}

	return uniqueClicks, nil
}

func (r *analyticsRepository) GetUniqueClicksByLinkID(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) (int64, error) {
	var uniqueClicks int64

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("link_clicks.link_id = ? AND links.user_id = ? AND links.deleted_at IS NULL", linkID, userID)

	db = applyAnalyticsDateFilter(db, "link_clicks.created_at", filter)

	if err := db.
		Select("COUNT(DISTINCT (link_clicks.ip_address, link_clicks.user_agent))").
		Scan(&uniqueClicks).Error; err != nil {
		return 0, err
	}

	return uniqueClicks, nil
}

func (r *analyticsRepository) GetTotalClicksByLinkID(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) (int64, error) {
	var totalClicks int64

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("link_clicks.link_id = ? AND links.user_id = ? AND links.deleted_at IS NULL", linkID, userID)

	db = applyAnalyticsDateFilter(db, "link_clicks.created_at", filter)

	if err := db.Count(&totalClicks).Error; err != nil {
		return 0, err
	}

	return totalClicks, nil
}

func (r *analyticsRepository) GetClicks(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) (analytics.GetClicksResult, error) {
	totalClicks, err := r.GetTotalClicksByLinkID(ctx, linkID, userID, filter)
	if err != nil {
		return analytics.GetClicksResult{}, err
	}

	uniqueClicks, err := r.GetUniqueClicksByLinkID(ctx, linkID, userID, filter)
	if err != nil {
		return analytics.GetClicksResult{}, err
	}

	return analytics.GetClicksResult{
		TotalClick:  totalClicks,
		UniqueClick: uniqueClicks,
	}, nil
}

func applyAnalyticsDateFilter(db *gorm.DB, column string, filter analytics.FilterQuery) *gorm.DB {
	if filter.FromDate == nil || filter.ToDate == nil {
		return db
	}

	return db.Where(column+" >= ? AND "+column+" < ?", filter.FromDate.UTC(), filter.ToDate.UTC())
}

func (r *analyticsRepository) GetEngagementByDevices(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) ([]analytics.GetTotalEngagementByDeviceResult, error) {
	var result []analytics.GetTotalEngagementByDeviceResult
	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("link_clicks.device, COUNT(*) AS clicks").
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("link_clicks.link_id = ? AND links.user_id = ? AND links.deleted_at IS NULL", linkID, userID)

	if filter.FromDate != nil && filter.ToDate != nil {
		db = db.Where("link_clicks.created_at >= ? AND link_clicks.created_at < ?", filter.FromDate.UTC(), filter.ToDate.UTC())
	}

	db = db.Group("link_clicks.device").
		Order("clicks DESC").
		Scan(&result)
	err := db.Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *analyticsRepository) GetEngagementByReferer(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) ([]analytics.GetTotalEngagementByRefererResult, error) {
	var result []analytics.GetTotalEngagementByRefererResult

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("COALESCE(NULLIF(referer, ''), 'direct') AS referer, COUNT(*) AS clicks").
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("link_clicks.link_id = ? AND links.user_id = ? AND links.deleted_at IS NULL", linkID, userID)

	if filter.FromDate != nil && filter.ToDate != nil {
		db = db.Where("link_clicks.created_at >= ? AND link_clicks.created_at < ?", filter.FromDate.UTC(), filter.ToDate.UTC())
	}

	if err := db.
		Group("COALESCE(NULLIF(referer, ''), 'direct')").
		Order("clicks DESC").
		Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *analyticsRepository) GetEngagementByOvertime(ctx context.Context, linkID, userID int64, filter analytics.FilterQuery) ([]analytics.GetTotalEngagementByOvertimeResult, error) {
	var result []analytics.GetTotalEngagementByOvertimeResult

	db := r.db.WithContext(ctx).
		Model(&entity.LinkClicksEntity{}).
		Select("DATE(link_clicks.created_at) AS date, COUNT(*) AS clicks").
		Joins("JOIN links ON links.id = link_clicks.link_id").
		Where("link_clicks.link_id = ? AND links.user_id = ? AND links.deleted_at IS NULL", linkID, userID)

	if filter.FromDate != nil && filter.ToDate != nil {
		db = db.Where("link_clicks.created_at >= ? AND link_clicks.created_at < ?", filter.FromDate.UTC(), filter.ToDate.UTC())
	}

	if err := db.
		Group("DATE(link_clicks.created_at)").
		Order("date DESC").
		Scan(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
