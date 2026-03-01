package account

import (
	"context"
	"errors"
	"shortivo/internal/http/response"

	"gorm.io/gorm"
)

type AccountService interface {
	GetProfile(ctx context.Context, userID int64) (GetProfile, error)
}

type accountService struct {
	accRepo AccountRepository
}

func NewService(accRepo AccountRepository) AccountService {
	return &accountService{accRepo: accRepo}
}

func (s *accountService) GetProfile(ctx context.Context, userID int64) (GetProfile, error) {
	profile, err := s.accRepo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GetProfile{}, response.BadRequestError("user not found")
		}

		return GetProfile{}, response.InternalServerError(err.Error())
	}

	return GetProfile{
		ID:        profile.ID,
		Username:  profile.Username,
		Email:     profile.Email,
		GoogleID:  profile.GoogleID,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}, nil
}
