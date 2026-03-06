package repo

import (
	"context"
	"shortivo/internal/app/subscription"
	"shortivo/internal/database/entity"

	"gorm.io/gorm"
)

type planRepository struct {
	db *gorm.DB
}

func NewPlanRespository(db *gorm.DB) subscription.PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) GetPlanByID(ctx context.Context, planID int64) (subscription.Plan, error) {
	var plan entity.PlanEntity
	if err := r.db.WithContext(ctx).Where("id = ?", planID).First(&plan).Error; err != nil {
		return subscription.Plan{}, err
	}

	return plan.ToDomain(), nil
}

func (r *planRepository) GetPlanByCode(ctx context.Context, code string) (subscription.Plan, error) {
	var plan entity.PlanEntity
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&plan).Error; err != nil {
		return subscription.Plan{}, err
	}

	return plan.ToDomain(), nil
}
