package auth

import (
	"errors"
	"shortivo/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWTConfig(cfg *config.Config) config.JwtConfig {
	return cfg.JWT
}

type JwtUser struct {
	UserId uint64
	Claims *jwt.MapClaims
}

type JwtManager interface {
	GenerateAccessToken(userId uint64) (string, error)
	ValidateAccessToken(accessToken string) (*JwtUser, error)
}

type jwtManager struct {
	cfg config.JwtConfig
}

func NewJWTManager(cfg config.JwtConfig) JwtManager {
	return &jwtManager{cfg: cfg}
}

func (jm *jwtManager) GenerateAccessToken(userId uint64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(jm.cfg.Expiration).Unix(),
	}

	return jm.generateJWT(claims, jm.cfg.AccessSecret)
}

func (jm *jwtManager) ValidateAccessToken(accessToken string) (*JwtUser, error) {
	claims, err := jm.validateJWT(accessToken, jm.cfg.AccessSecret)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	if !jm.checkClaims(claims, "userId", "exp") {
		return nil, errors.New("invalid token")
	}

	uidFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, errors.New("invalid userId claim")
	}

	return &JwtUser{
		UserId: uint64(uidFloat),
		Claims: &claims,
	}, nil
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
		return nil, err
	}

	claims, _ := j.Claims.(jwt.MapClaims)

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return nil, errors.New("token expired")
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
