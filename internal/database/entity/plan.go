package entity

type PlanEntity struct {
	ID                  int64  `gorm:"primaryKey"`
	Title               string `gorm:"index"`
	MonthlyLimit        int16
	AnalyticsEnabled    bool
	CustomDomainEnabled bool
}
