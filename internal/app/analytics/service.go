package analytics

type AnalyticsService interface{}

type analyticsService struct{}

func NewService() AnalyticsService {
	return &analyticsService{}
}
