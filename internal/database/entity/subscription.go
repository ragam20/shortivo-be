package entity

import "time"

type SubscriptionStatus string

const (
	SubscriptionActive    SubscriptionStatus = "active"
	SubscriptionExpired   SubscriptionStatus = "expired"
	SubscriptionCancelled SubscriptionStatus = "cancelled"
)

type SubscriptionEntity struct {
	ID        int64              `gorm:"primaryKey;autoIncrement:false"`
	UserID    int64              `gorm:"not null;index:idx_subscription_user_status"`
	User      UserEntity         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	PlanID    int64              `gorm:"not null;index"`
	Plan      PlanEntity         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:PlanID;references:ID"`
	PeriodEnd time.Time          `gorm:"not null;index"`
	Status    SubscriptionStatus `gorm:"type:varchar(20);not null;default:'active';index:idx_subscription_user_status;check:status IN ('active','expired','cancelled')"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
