package link

type LinkService interface{}

type linkService struct{}

func NewService() LinkService {
	return &linkService{}
}
