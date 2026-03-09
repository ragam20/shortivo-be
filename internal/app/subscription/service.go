package subscription

import (
	"context"
	"errors"

	"shortivo/internal/http/response"

	"gorm.io/gorm"
)

type SubscriptionService interface {
	CancelSubscription(ctx context.Context, cmd CancelSubscriptionCMD) error
	GetCurrentPlan(ctx context.Context, userID int64) (Plan, error)
	GetCurrentSubscription(ctx context.Context, userID int64) (GetCurrSubscriptionResult, error)
}

type subscriptionService struct {
	subsRepo SubscriptionRepository
	planRepo PlanRepository
}

func NewService(subRepo SubscriptionRepository, planRepo PlanRepository) SubscriptionService {
	return &subscriptionService{
		subsRepo: subRepo,
		planRepo: planRepo,
	}
}

func (s *subscriptionService) GetCurrentSubscription(ctx context.Context, userID int64) (GetCurrSubscriptionResult, error) {
	var result GetCurrSubscriptionResult

	sub, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, response.InternalServerError("error GetCurrentSubscription")
		}
	}

	if sub == nil {
		plan, _ := s.planRepo.GetPlanByCode(ctx, "free")

		result.Plan = PlanResult{
			PlanID:              plan.ID,
			Title:               plan.Title,
			Code:                plan.Code,
			MonthlyLimit:        plan.MonthlyLimit,
			AnalyticsEnabled:    plan.AnalyticsEnabled,
			CustomDomainEnabled: plan.CustomDomainEnabled,
		}

		return result, nil
	}

	plan, _ := s.planRepo.GetPlanByID(ctx, sub.PlanID)
	result.Plan = PlanResult{
		PlanID:              plan.ID,
		Title:               plan.Title,
		Code:                plan.Code,
		MonthlyLimit:        plan.MonthlyLimit,
		AnalyticsEnabled:    plan.AnalyticsEnabled,
		CustomDomainEnabled: plan.CustomDomainEnabled,
	}
	result.Subscription = &SubscriptionResult{
		ID:        sub.ID,
		PeriodEnd: sub.PeriodEnd,
		Status:    sub.Status,
	}

	return result, nil
}

func (s *subscriptionService) GetCurrentPlan(ctx context.Context, userID int64) (Plan, error) {
	var plan Plan

	subs, err := s.subsRepo.GetActiveSubscriptionByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			plan, err = s.planRepo.GetPlanByCode(ctx, "free")
			if err != nil {
				return Plan{}, response.InternalServerError("erro GetPlanByID")
			}
		} else {
			return Plan{}, response.InternalServerError("error GetActiveSubscriptionByUser")
		}
	} else {
		plan, err = s.planRepo.GetPlanByID(ctx, subs.PlanID)
		if err != nil {
			return Plan{}, err
		}
	}

	return plan, nil
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
