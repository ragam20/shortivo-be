package auth

type AuthService interface{}

type authService struct{}

func NewService() AuthService {
	return &authService{}
}
