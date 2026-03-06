package subscription

import "context"

type PlanRepository interface {
	GetPlanByID(ctx context.Context, planID int64) (Plan, error)
	GetPlanByCode(ctx context.Context, code string) (Plan, error)
}
