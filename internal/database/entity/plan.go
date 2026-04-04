package entity

import (
	"time"

	"shortivo/internal/app/subscription"
)

type PlanEntity struct {
	ID                  int64  `gorm:"primaryKey"`
	Code                string `gorm:"index"`
	Title               string `gorm:"size:64;not null;uniqueIndex"`
	MonthlyLimit        int32  `gorm:"not null;default:0"`
	AnalyticsEnabled    bool   `gorm:"not null;default:false"`
	CustomDomainEnabled bool   `gorm:"not null;default:false"`

	UpdatedAt time.Time
	CreatedAt time.Time
}

func (PlanEntity) TableName() string {
	return "plans"
}

func NewPlanEntity(domain subscription.Plan) PlanEntity {
	return PlanEntity{
		ID:                  domain.ID,
		Title:               domain.Title,
		Code:                domain.Code,
		MonthlyLimit:        domain.MonthlyLimit,
		AnalyticsEnabled:    domain.AnalyticsEnabled,
		CustomDomainEnabled: domain.CustomDomainEnabled,
		UpdatedAt:           domain.UpdatedAt,
		CreatedAt:           domain.CreatedAt,
	}
}

func (e PlanEntity) ToDomain() subscription.Plan {
	return subscription.Plan{
		ID:                  e.ID,
		Title:               e.Title,
		Code:                e.Code,
		MonthlyLimit:        e.MonthlyLimit,
		AnalyticsEnabled:    e.AnalyticsEnabled,
		CustomDomainEnabled: e.CustomDomainEnabled,
		UpdatedAt:           e.UpdatedAt,
		CreatedAt:           e.CreatedAt,
	}
}
