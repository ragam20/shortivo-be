package server

import (
	"net/http"
	"time"

	"shortivo/internal/config"
)

const (
	googleOAuthStateCookieName = "googleOAuthState"
	refreshTokenCookie         = "refreshToken"
)

type CookieManager struct {
	appCfg *config.AppConfig
}

func NewCookieManager(appCfg *config.AppConfig) *CookieManager {
	return &CookieManager{appCfg: appCfg}
}

func (m *CookieManager) SetRefreshToken(refToken string, expiredAt time.Time) *http.Cookie {
	maxAge := max(int(time.Until(expiredAt).Seconds()), 0)

	return &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refToken,
		Path:     "/",
		Secure:   m.isSecure(),
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (m *CookieManager) RevokeRefreshToken() *http.Cookie {
	return &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/",
		Secure:   m.isSecure(),
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (m *CookieManager) SetGoogleOAuthState(state string) *http.Cookie {
	return &http.Cookie{
		Name:     googleOAuthStateCookieName,
		Value:    state,
		Path:     "/",
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
		Path:     "/",
		Secure:   m.isSecure(),
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
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
