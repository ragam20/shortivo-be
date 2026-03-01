package repo

import (
	"shortivo/internal/app/analytics"

	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) analytics.AnalyticsRepository {
	return &analyticsRepository{db: db}
}
