package subscription

import (
	"context"
	"errors"

	"shortivo/internal/http/response"

	"gorm.io/gorm"
)

type SubscriptionService interface {
	CancelSubscription(ctx context.Context, cmd CancelSubscriptionCMD) error
}

type subscriptionService struct {
	subsRepo SubscriptionRepository
}

func NewService(subRepo SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subsRepo: subRepo,
	}
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, cmd CancelSubscriptionCMD) error {
	if cmd.UserID <= 0 {
		return response.Unauthorized()
	}

	if cmd.SubsID <= 0 {
		return response.BadRequestError("subsID is required")
	}

	if err := s.subsRepo.ChangeStatus(ctx, cmd.UserID, cmd.SubsID, SubscriptionCancelled); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFoundError("subsription not found")
		}
	}

	return nil
}
