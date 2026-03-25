package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"shortivo/internal/app/auth/oauth"
	"shortivo/internal/config"
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
	GoogleSignInURL() (string, string, error)
	SignInWithGoogle(ctx context.Context, req request.GoogleOAuthCallbackRequest, expectedState string, userAgent string, ipAddress string) (AuthTokens, error)
	SignOut(ctx context.Context, userID int64)
	RefreshToken(ctx context.Context, token string, userAgent string, ipAddress string) (AuthTokens, error)
}

type authService struct {
	idGen       idgen.IDGenerator
	authRepo    AuthRepository
	jwtManager  JwtManager
	googleOAuth oauth.GoogleOAuthProvider
}

func NewService(
	idGen idgen.IDGenerator,
	authRepo AuthRepository,
	jwtManager JwtManager,
	googleOAuth oauth.GoogleOAuthProvider,
	_ *config.GoogleOAuthConfig,
) AuthService {
	return &authService{
		idGen:       idGen,
		authRepo:    authRepo,
		jwtManager:  jwtManager,
		googleOAuth: googleOAuth,
	}
}

func (s *authService) SignIn(ctx context.Context, req request.SignInRequest, userAgent string, ipAddress string) (AuthTokens, error) {
	if strings.TrimSpace(req.Email) == "" {
		return AuthTokens{}, response.BadRequestError("email is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return AuthTokens{}, response.BadRequestError("password is required")
	}

	auth, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AuthTokens{}, response.NotFoundError("user not found")
		}

		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	if auth.Password == nil || strings.TrimSpace(*auth.Password) == "" {
		return AuthTokens{}, response.ForbiddenError("this account uses google sign in")
	}

	if utils.ValidatePassword(*auth.Password, req.Password) {
		return AuthTokens{}, response.BadRequestError("wrong password")
	}

	return s.issueTokens(ctx, s.authRepo, auth.ID, userAgent, ipAddress)
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

	var tokens AuthTokens

	if err := s.authRepo.WithTransaction(ctx, func(repo AuthRepository) error {
		if err := repo.SignUp(ctx, auth); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return response.BadRequestError("user already exists")
			}

			return response.InternalServerError(fmt.Sprintf("error signup: %s", err.Error()))
		}

		issuedTokens, err := s.issueTokens(ctx, repo, auth.ID, userAgent, ipAddress)
		if err != nil {
			return err
		}
		tokens = issuedTokens

		return nil
	}); err != nil {
		return AuthTokens{}, err
	}

	return tokens, nil
}

func (s *authService) SignOut(ctx context.Context, userID int64) {
	s.authRepo.RevokeTokenByUser(ctx, userID)
}

func (s *authService) GoogleSignInURL() (string, string, error) {
	if !s.googleOAuth.Enabled() {
		return "", "", response.NewAppErr(http.StatusServiceUnavailable, "google oauth is not configured")
	}

	state, err := s.generateSecureToken()
	if err != nil {
		return "", "", response.InternalServerError(err.Error())
	}

	return s.googleOAuth.GetURL(state), state, nil
}

func (s *authService) SignInWithGoogle(
	ctx context.Context,
	req request.GoogleOAuthCallbackRequest,
	expectedState string,
	userAgent string,
	ipAddress string,
) (AuthTokens, error) {
	if !s.googleOAuth.Enabled() {
		return AuthTokens{}, response.NewAppErr(http.StatusServiceUnavailable, "google oauth is not configured")
	}

	if strings.TrimSpace(expectedState) == "" || req.State != expectedState {
		return AuthTokens{}, response.Unauthorized()
	}

	gUser, err := s.googleOAuth.GetUser(ctx, req.Code)
	if err != nil {
		return AuthTokens{}, response.BadRequestError("invalid google oauth callback")
	}

	if !gUser.EmailVerified {
		return AuthTokens{}, response.ForbiddenError("google email must be verified")
	}

	if strings.TrimSpace(gUser.Sub) == "" || strings.TrimSpace(gUser.Email) == "" {
		return AuthTokens{}, response.BadRequestError("invalid google user data")
	}

	var tokens AuthTokens
	if err := s.authRepo.WithTransaction(ctx, func(repo AuthRepository) error {
		googleAuth, err := s.findOrCreateGoogleUser(ctx, repo, gUser)
		if err != nil {
			return err
		}

		issuedTokens, err := s.issueTokens(ctx, repo, googleAuth.ID, userAgent, ipAddress)
		if err != nil {
			return err
		}

		tokens = issuedTokens
		return nil
	}); err != nil {
		return AuthTokens{}, err
	}

	return tokens, nil
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

	return s.issueTokens(ctx, s.authRepo, stored.UserID, userAgent, ipAddress)

}

func (s *authService) findOrCreateGoogleUser(ctx context.Context, repo AuthRepository, gUser *oauth.GoogleUser) (Auth, error) {
	auth, err := repo.FindByGoogleID(ctx, gUser.Sub)
	if err == nil {
		return auth, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return Auth{}, response.InternalServerError(err.Error())
	}

	auth, err = repo.FindByEmail(ctx, gUser.Email)
	if err == nil {
		if auth.GoogleID != nil && strings.TrimSpace(*auth.GoogleID) != "" {
			if *auth.GoogleID != gUser.Sub {
				return Auth{}, response.ForbiddenError("google account does not match existing user")
			}
			return auth, nil
		}

		if err := repo.UpdateGoogleID(ctx, auth.ID, gUser.Sub); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return s.findGoogleUserAfterConflict(ctx, repo, gUser)
			}
			return Auth{}, response.InternalServerError(err.Error())
		}

		auth.GoogleID = &gUser.Sub
		return auth, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return Auth{}, response.InternalServerError(err.Error())
	}

	newAuth := Auth{
		ID:       s.idGen.NewID(),
		Email:    gUser.Email,
		Username: buildGoogleUsername(gUser.Name, gUser.Email, gUser.Sub),
		GoogleID: &gUser.Sub,
	}

	if err := repo.SignUp(ctx, newAuth); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.findGoogleUserAfterConflict(ctx, repo, gUser)
		}
		return Auth{}, response.InternalServerError(fmt.Sprintf("error signup: %s", err.Error()))
	}

	return newAuth, nil
}

func (s *authService) findGoogleUserAfterConflict(ctx context.Context, repo AuthRepository, gUser *oauth.GoogleUser) (Auth, error) {
	auth, err := repo.FindByGoogleID(ctx, gUser.Sub)
	if err == nil {
		return auth, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Auth{}, response.InternalServerError(err.Error())
	}

	auth, err = repo.FindByEmail(ctx, gUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Auth{}, response.BadRequestError("user already exists")
		}
		return Auth{}, response.InternalServerError(err.Error())
	}

	if auth.GoogleID != nil && *auth.GoogleID == gUser.Sub {
		return auth, nil
	}

	return Auth{}, response.ForbiddenError("google account does not match existing user")
}

func (s *authService) issueTokens(ctx context.Context, repo AuthRepository, userID int64, userAgent string, ipAddress string) (AuthTokens, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return AuthTokens{}, err
	}

	secureToken, err := s.generateSecureToken()
	if err != nil {
		return AuthTokens{}, response.InternalServerError(err.Error())
	}

	refreshToken := s.jwtManager.NewRefreshToken(s.idGen.NewID(), userID, secureToken, userAgent, ipAddress)

	if err := repo.SaveRefreshToken(ctx, refreshToken); err != nil {
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

func buildGoogleUsername(name string, email string, googleID string) string {
	for _, candidate := range []string{name, strings.Split(email, "@")[0]} {
		username := sanitizeUsername(candidate)
		if len(username) >= 5 {
			return truncateUsername(username)
		}
	}

	return truncateUsername("user_" + utils.HashString(email + googleID)[:10])
}

func sanitizeUsername(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return ""
	}

	var builder strings.Builder
	lastWasSeparator := false

	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
			lastWasSeparator = false
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastWasSeparator = false
		case !lastWasSeparator && builder.Len() > 0:
			builder.WriteByte('_')
			lastWasSeparator = true
		}
	}

	return strings.Trim(builder.String(), "_")
}

func truncateUsername(username string) string {
	if len(username) <= 50 {
		return username
	}
	return username[:50]
}
