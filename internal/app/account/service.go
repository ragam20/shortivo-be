package account

type AccountService interface{}

type accountService struct{}

func NewService() AccountService {
	return &accountService{}
}
