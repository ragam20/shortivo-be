package auth

import (
	"shortivo/internal/config"
	"shortivo/internal/http/response"
	"shortivo/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWTConfig(cfg *config.Config) config.JwtConfig {
	return cfg.JWT
}

type JwtUser struct {
	UserId int64
	Claims *jwt.MapClaims
}

type JwtManager interface {
	GenerateAccessToken(userId int64) (string, error)
	ValidateAccessToken(accessToken string) (*JwtUser, error)
	NewRefreshToken(id int64, userID int64, token string, userAgent string, ipAddr string) RefreshToken
}

type jwtManager struct {
	cfg config.JwtConfig
}

func NewJWTManager(cfg config.JwtConfig) JwtManager {
	return &jwtManager{cfg: cfg}
}

func (jm *jwtManager) GenerateAccessToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(jm.cfg.Expiration).Unix(),
	}

	return jm.generateJWT(claims, jm.cfg.AccessSecret)
}

func (jm *jwtManager) ValidateAccessToken(accessToken string) (*JwtUser, error) {
	claims, err := jm.validateJWT(accessToken, jm.cfg.AccessSecret)
	if err != nil {
		return nil, response.Unauthorized()
	}

	if !jm.checkClaims(claims, "userId", "exp") {
		return nil, response.BadRequestError("invalid jwt")
	}

	uidFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, response.BadRequestError("invalid userId claim")
	}

	return &JwtUser{
		UserId: int64(uidFloat),
		Claims: &claims,
	}, nil
}

func (jm *jwtManager) NewRefreshToken(id int64, userID int64, token string, userAgent string, ipAddr string) RefreshToken {
	refToken := RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: utils.HashString(token),
		UserAgent: userAgent,
		IpAddress: ipAddr,
		ExpiredAt: time.Now().Add(jm.cfg.RefreshExp),
	}

	return refToken
}

func (jm *jwtManager) generateJWT(jwtClaims jwt.MapClaims, secret string) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtClaims,
	)
	return t.SignedString([]byte(secret))
}

func (jm *jwtManager) validateJWT(currToken string, secret string) (jwt.MapClaims, error) {
	j, err := jwt.Parse(currToken, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !j.Valid {
		return nil, response.BadRequestError("invalid jwt")
	}

	claims, _ := j.Claims.(jwt.MapClaims)

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return nil, response.BadRequestError("token expired")
	}

	return claims, nil
}

func (jm *jwtManager) checkClaims(claims jwt.MapClaims, keys ...string) bool {
	for _, key := range keys {
		if _, ok := claims[key]; !ok {
			return false
		}
	}
	return true
}
