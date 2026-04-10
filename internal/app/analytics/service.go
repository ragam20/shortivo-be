package analytics

import (
	"context"
	"errors"

	"shortivo/internal/app/subscription"
	"shortivo/internal/http/response"

	"gorm.io/gorm"
)

type AnalyticsService interface {
	GetTotalCreatedLink(ctx context.Context, userID int64, filter FilterQuery) (int64, error)
	GetTotalClicks(ctx context.Context, userID int64, filter FilterQuery) (GetClicksResult, error)
	GetTotalClicksByLinkID(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) (GetClicksResult, error)
	GetEngagementByDevices(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) ([]GetTotalEngagementByDeviceResult, error)
	GetEngagementByReferer(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) ([]GetTotalEngagementByRefererResult, error)
	GetTotalEngagementByOvertime(ctx context.Context, query GetTotalEngagementByOvertimeQuery, filter FilterQuery) ([]GetTotalEngagementByOvertimeResult, error)
}

type analyticsService struct {
	analyticsRepo AnalyticsRepository
	subsRepo      subscription.SubscriptionRepository
}

func NewService(analitycsRepo AnalyticsRepository, subsRepo subscription.SubscriptionRepository) AnalyticsService {
	return &analyticsService{analyticsRepo: analitycsRepo, subsRepo: subsRepo}
}

func (s *analyticsService) GetTotalCreatedLink(ctx context.Context, userID int64, filter FilterQuery) (int64, error) {
	total, err := s.analyticsRepo.GetTotalCreatedLink(ctx, userID, filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, response.BadRequestError("links not found")
		}
		return 0, response.InternalServerError(err.Error())
	}

	return total, nil
}

func (s *analyticsService) GetTotalClicks(ctx context.Context, userID int64, filter FilterQuery) (GetClicksResult, error) {
	clickCount, err := s.analyticsRepo.GetTotalClicks(ctx, userID, filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GetClicksResult{}, response.BadRequestError("click event not found")
		}
		return GetClicksResult{}, response.InternalServerError(err.Error())
	}

	uniqueCount, err := s.analyticsRepo.GetUniqueClicks(ctx, userID, filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GetClicksResult{}, response.BadRequestError("unique click event not found")
		}
		return GetClicksResult{}, response.InternalServerError(err.Error())
	}

	return GetClicksResult{
		TotalClick:  clickCount,
		UniqueClick: uniqueCount,
	}, nil
}

func (s *analyticsService) GetTotalClicksByLinkID(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) (GetClicksResult, error) {
	if query.UserID <= 0 {
		return GetClicksResult{}, response.Unauthorized()
	}

	if query.LinkID <= 0 {
		return GetClicksResult{}, response.BadRequestError("linkID is required")
	}

	_, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GetClicksResult{}, response.ForbiddenError("pro subscription required")
		}
		return GetClicksResult{}, response.InternalServerError("error GetActiveSubscriptionByUser")
	}

	result, err := s.analyticsRepo.GetClicks(ctx, query.LinkID, query.UserID, filter)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return GetClicksResult{}, response.InternalServerError(err.Error())
		}
	}

	return result, nil
}

func (s *analyticsService) GetEngagementByDevices(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) ([]GetTotalEngagementByDeviceResult, error) {
	if query.UserID <= 0 {
		return []GetTotalEngagementByDeviceResult{}, response.Unauthorized()
	}

	if query.LinkID <= 0 {
		return []GetTotalEngagementByDeviceResult{}, response.BadRequestError("link id is required")
	}

	_, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByDeviceResult{}, response.ForbiddenError("pro subscription required")
		}
		return []GetTotalEngagementByDeviceResult{}, response.InternalServerError("error GetActiveSubscriptionByUser")
	}

	result, err := s.analyticsRepo.GetEngagementByDevices(ctx, query.LinkID, query.UserID, filter)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByDeviceResult{}, response.InternalServerError("error GetTotalEngagementByDeviceResult")
		}
	}

	return result, nil
}

func (s *analyticsService) GetEngagementByReferer(ctx context.Context, query GetTotalClicksQuery, filter FilterQuery) ([]GetTotalEngagementByRefererResult, error) {
	if query.UserID <= 0 {
		return []GetTotalEngagementByRefererResult{}, response.Unauthorized()
	}

	if query.LinkID <= 0 {
		return []GetTotalEngagementByRefererResult{}, response.BadRequestError("link id is required")
	}

	_, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByRefererResult{}, response.ForbiddenError("pro subscription required")
		}
		return []GetTotalEngagementByRefererResult{}, response.InternalServerError("error GetActiveSubscriptionByUser")
	}

	result, err := s.analyticsRepo.GetEngagementByReferer(ctx, query.LinkID, query.UserID, filter)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByRefererResult{}, response.InternalServerError("error GetTotalEngagementByRefererResult")
		}
	}

	return result, nil
}

func (s *analyticsService) GetTotalEngagementByOvertime(ctx context.Context, query GetTotalEngagementByOvertimeQuery, filter FilterQuery) ([]GetTotalEngagementByOvertimeResult, error) {
	if query.UserID <= 0 {
		return []GetTotalEngagementByOvertimeResult{}, response.Unauthorized()
	}

	if query.LinkID <= 0 {
		return []GetTotalEngagementByOvertimeResult{}, response.BadRequestError("link id is required")
	}

	_, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByOvertimeResult{}, response.ForbiddenError("pro subscription required")
		}
		return []GetTotalEngagementByOvertimeResult{}, response.InternalServerError("error GetActiveSubscriptionByUser")
	}

	result, err := s.analyticsRepo.GetEngagementByOvertime(ctx, query.LinkID, query.UserID, filter)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetTotalEngagementByOvertimeResult{}, response.InternalServerError("error GetTotalEngagementByOvertime")
		}
	}

	return result, nil
}
