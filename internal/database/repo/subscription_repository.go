package repo

import (
	"context"
	"shortivo/internal/app/subscription"
	"shortivo/internal/database/entity"
	"time"

	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) GetSubscriptionByUser(ctx context.Context, userID int64) (subscription.Subscription, error) {
	var sub entity.SubscriptionEntity
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND period_end > ?", userID, subscription.SubscriptionActive, time.Now()).
		Order("period_end DESC").
		First(&sub).Error; err != nil {
		return subscription.Subscription{}, err
	}

	return sub.ToDomain(), nil
}
