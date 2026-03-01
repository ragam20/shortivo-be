package entity

import "time"

type PlanEntity struct {
	ID                  int64  `gorm:"primaryKey"`
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
