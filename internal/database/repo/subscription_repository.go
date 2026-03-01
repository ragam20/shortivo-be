package repo

import (
	"shortivo/internal/app/subscription"

	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}
