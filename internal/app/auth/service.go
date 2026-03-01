package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"
	"shortivo/internal/idgen"
	"shortivo/internal/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	SignIn(ctx context.Context, req request.SignInRequest, userAgent string, ipAddress string) (AuthTokens, error)
	SignUp(ctx context.Context, req request.SignUpRequest, userAgent string, ipAddress string) (AuthTokens, error)
	SignOut(ctx context.Context, userID int64)
	RefreshToken(ctx context.Context, token string, userAgent string, ipAddress string) (AuthTokens, error)
}

type authService struct {
	idGen      idgen.IDGenerator
	authRepo   AuthRepository
	jwtManager JwtManager
}

func NewService(idGen idgen.IDGenerator, authRepo AuthRepository, jwtManager JwtManager) AuthService {
	return &authService{
		idGen:      idGen,
		authRepo:   authRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) SignIn(ctx context.Context, req request.SignInRequest, userAgent string, ipAddress string) (AuthTokens, error) {
	if strings.TrimSpace(req.Email) == "" {
		return AuthTokens{}, response.BadRequestError("email is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return AuthTokens{}, response.BadRequestError("password is required")
	}

	auth, err := s.authRepo.SignIn(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthTokens{}, response.NotFoundError("user not found")
		}

		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	if utils.ValidatePassword(*auth.Password, req.Password) {
		return AuthTokens{}, response.BadRequestError("wrong password")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(auth.ID)
	if err != nil {
		return AuthTokens{}, err
	}

	secureToken, err := s.generateSecureToken()
	if err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	refreshToken := s.jwtManager.NewRefreshToken(s.idGen.NewID(), auth.ID, secureToken, userAgent, ipAddress)

	if err := s.authRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: secureToken,
		RefreshExp:   refreshToken.ExpiredAt,
	}, nil
}

func (s *authService) SignUp(ctx context.Context, req request.SignUpRequest, userAgent string, ipAddress string) (AuthTokens, error) {
	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	auth := Auth{
		ID:       s.idGen.NewID(),
		Email:    req.Email,
		Username: req.Username,
		Password: &hashPass,
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(auth.ID)
	if err != nil {
		return AuthTokens{}, err
	}

	secureToken, err := s.generateSecureToken()
	if err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	var refreshToken RefreshToken

	if err := s.authRepo.WithTransaction(ctx, func(repo AuthRepository) error {
		if err := s.authRepo.SignUp(ctx, auth); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return response.BadRequestError("user already exists")
			}

			return response.InternalServerError(fmt.Sprintf("error signup: %s", err.Error()))
		}

		refreshToken = s.jwtManager.NewRefreshToken(s.idGen.NewID(), auth.ID, secureToken, userAgent, ipAddress)

		if err := s.authRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
			return response.InternalServerError(err.Error())
		}

		return nil
	}); err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: secureToken,
		RefreshExp:   refreshToken.ExpiredAt,
	}, nil
}

func (s *authService) SignOut(ctx context.Context, userID int64) {
	s.authRepo.RevokeTokenByUser(ctx, userID)
}

func (s *authService) RefreshToken(ctx context.Context, token string, userAgent string, ipAddress string) (AuthTokens, error) {
	hashToken := utils.HashString(token)

	stored, err := s.authRepo.FindByTokenHash(ctx, hashToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthTokens{}, response.Unauthorized()
		}
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	if stored.Revoked || time.Now().After(stored.ExpiredAt) {
		return AuthTokens{}, response.Unauthorized()
	}

	if err := s.authRepo.DeleteToken(ctx, hashToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthTokens{}, response.Unauthorized()
		}
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(stored.UserID)
	if err != nil {
		return AuthTokens{}, err
	}

	secureToken, err := s.generateSecureToken()
	if err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	refreshToken := s.jwtManager.NewRefreshToken(s.idGen.NewID(), stored.UserID, secureToken, userAgent, ipAddress)

	if err := s.authRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: secureToken,
		RefreshExp:   refreshToken.ExpiredAt,
	}, nil

}

func (s *authService) generateSecureToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
