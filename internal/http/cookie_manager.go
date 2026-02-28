package server

import (
	"net/http"
	"shortivo/internal/config"
	"time"
)

type CookieManager struct {
	appCfg *config.AppConfig
}

func NewCookieManager(appCfg *config.AppConfig) *CookieManager {
	return &CookieManager{appCfg: appCfg}
}

func (m *CookieManager) SetRefreshToken(refToken string, expiredAt time.Time) *http.Cookie {
	secure := false
	if m.appCfg.Env != "development" {
		secure = true
	}

	maxAge := max(int(time.Until(expiredAt).Seconds()), 0)

	return &http.Cookie{
		Name:     "refreshToken",
		Value:    refToken,
		Path:     "/refresh-token",
		Secure:   secure,
		MaxAge:   maxAge,
		HttpOnly: true,
	}
}

func (m *CookieManager) RevokeRefreshToken() *http.Cookie {
	return &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/refresh-token",
		MaxAge:   0,
		HttpOnly: true,
	}
}
