package response

import "time"

type PlanResponse struct {
	ID                  int64  `json:"id"`
	Code                string `json:"code"`
	Title               string `json:"title"`
	MonthlyLimit        int32  `json:"monthlyLimit"`
	AnalyticsEnabled    bool   `json:"analyticsEnabled"`
	CustomDomainEnabled bool   `json:"customDomainEnabled"`
}

type SubscriptionResponse struct {
	ID        int64     `json:"id"`
	PeriodEnd time.Time `json:"periodEnd"`
	Status    string    `json:"status"`
}

type GetCurrentSubscriptionResponse struct {
	Plan         PlanResponse          `json:"plan"`
	Subscription *SubscriptionResponse `json:"subscription"`
}
