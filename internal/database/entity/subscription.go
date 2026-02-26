package entity

import "time"

type SubscriptionStatus string

const (
	SubscriptionActive    SubscriptionStatus = "active"
	SubscriptionExpired   SubscriptionStatus = "expired"
	SubscriptionCancelled SubscriptionStatus = "cancelled"
)

type SubscriptionEntity struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	PlanID    int64
	PeriodEnd time.Time
	Status    SubscriptionStatus `gorm:"type:varchar(20);not null;default:'active'"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
