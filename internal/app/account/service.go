package account

import (
	"context"
	"errors"

	"shortivo/internal/app/subscription"
	"shortivo/internal/http/response"

	"gorm.io/gorm"
)

type AccountService interface {
	GetAccount(ctx context.Context, userID int64) (UserPlan, error)
}

type accountService struct {
	accRepo  AccountRepository
	subsRepo subscription.SubscriptionRepository
}

func NewService(accRepo AccountRepository, subsRepo subscription.SubscriptionRepository) AccountService {
	return &accountService{
		accRepo:  accRepo,
		subsRepo: subsRepo,
	}
}

func (s *accountService) GetAccount(ctx context.Context, userID int64) (UserPlan, error) {
	profile, err := s.accRepo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserPlan{}, response.BadRequestError("user not found")
		}

		return UserPlan{}, response.InternalServerError(err.Error())
	}

	account := GetAccount{
		ID:        profile.ID,
		Username:  profile.Username,
		Email:     profile.Email,
		GoogleID:  profile.GoogleID,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	var plan *Plan
	pl, err := s.subsRepo.GetCurrentPlanByUserID(ctx, userID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return UserPlan{}, response.InternalServerError(err.Error())
		}
	} else {
		plan = &Plan{
			ID:                  pl.ID,
			Title:               pl.Title,
			Code:                pl.Code,
			MonthlyLimit:        pl.MonthlyLimit,
			AnalyticsEnabled:    pl.AnalyticsEnabled,
			CustomDomainEnabled: pl.CustomDomainEnabled,
		}
	}

	return UserPlan{
		User: account,
		Plan: plan,
	}, nil
}
