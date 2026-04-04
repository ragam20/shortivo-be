package repo

import (
	"context"
	"time"

	"shortivo/internal/app/subscription"
	"shortivo/internal/database/entity"

	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) GetActiveSubscriptionByUser(ctx context.Context, userID int64) (*subscription.Subscription, error) {
	var sub entity.SubscriptionEntity
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND period_end > ?", userID, subscription.SubscriptionActive, time.Now()).
		Order("period_end DESC").
		First(&sub).Error; err != nil {
		return nil, err
	}

	domain := sub.ToDomain()
	return &domain, nil
}

func (r *subscriptionRepository) ChangeStatus(ctx context.Context, userID int64, subsID int64, sub subscription.SubscriptionStatus) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.SubscriptionEntity{}).
		Where("id = ? AND user_id = ?", subsID, userID).
		UpdateColumn("status", sub).Error; err != nil {
		return err
	}

	return nil
}

func (r *subscriptionRepository) CreateSubscription(ctx context.Context, userID int64, subs subscription.Subscription) error {
	e := entity.NewSubscriptionEntity(subs)
	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return err
	}

	return nil
}

func (r *subscriptionRepository) GetCurrentPlanByUserID(ctx context.Context, userID int64) (subscription.Plan, error) {
	var sub entity.SubscriptionEntity
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND period_end > ?", userID, subscription.SubscriptionActive, time.Now()).
		Preload("Plan").
		Order("period_end DESC").
		First(&sub).Error; err != nil {
		return subscription.Plan{}, err
	}

	return sub.Plan.ToDomain(), nil
}
