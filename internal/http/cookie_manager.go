package server

import (
	"net/http"
	"shortivo/internal/config"
	"time"
)

const googleOAuthStateCookieName = "googleOAuthState"

type CookieManager struct {
	appCfg *config.AppConfig
}

func NewCookieManager(appCfg *config.AppConfig) *CookieManager {
	return &CookieManager{appCfg: appCfg}
}

func (m *CookieManager) SetRefreshToken(refToken string, expiredAt time.Time) *http.Cookie {
	maxAge := max(int(time.Until(expiredAt).Seconds()), 0)

	return &http.Cookie{
		Name:     "refreshToken",
		Value:    refToken,
		Path:     "/api/auth/refresh-token",
		Secure:   m.isSecure(),
		MaxAge:   maxAge,
		HttpOnly: true,
	}
}

func (m *CookieManager) RevokeRefreshToken() *http.Cookie {
	return &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/api/auth/refresh-token",
		Secure:   m.isSecure(),
		MaxAge:   0,
		HttpOnly: true,
	}
}

func (m *CookieManager) SetGoogleOAuthState(state string) *http.Cookie {
	return &http.Cookie{
		Name:     googleOAuthStateCookieName,
		Value:    state,
		Path:     "/api/auth/google/redirect",
		Secure:   m.isSecure(),
		MaxAge:   int((10 * time.Minute).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (m *CookieManager) RevokeGoogleOAuthState() *http.Cookie {
	return &http.Cookie{
		Name:     googleOAuthStateCookieName,
		Value:    "",
		Path:     "/api/auth/google/redirect",
		Secure:   m.isSecure(),
		MaxAge:   0,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func GoogleOAuthStateCookieName() string {
	return googleOAuthStateCookieName
}

func (m *CookieManager) isSecure() bool {
	return m.appCfg.Env != "development"
}
