package subscription

type SubscriptionService interface{}

type subscriptionService struct{}

func NewService() SubscriptionService {
	return &subscriptionService{}
}
