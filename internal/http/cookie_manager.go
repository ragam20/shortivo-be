package server

import (
	"net/http"
	"time"

	"shortivo/internal/config"
)

const (
	googleOAuthStateCookieName = "googleOAuthState"
	refreshTokenCookie         = "refreshToken"
	accessTokenCookie          = "authSession"
)

type CookieManager struct {
	appCfg *config.AppConfig
}

func NewCookieManager(appCfg *config.AppConfig) *CookieManager {
	return &CookieManager{appCfg: appCfg}
}

func (m *CookieManager) SetAccessToken(accessToken string, expiredAt time.Time) *http.Cookie {
	maxAge := max(int(time.Until(expiredAt).Seconds()), 0)

	return &http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		Secure:   m.isSecure(),
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: m.isSameSite(),
	}
}

func (m *CookieManager) RevokeAccessToken() *http.Cookie {
	return &http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		Secure:   m.isSecure(),
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		SameSite: m.isSameSite(),
	}
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
		SameSite: m.isSameSite(),
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
		SameSite: m.isSameSite(),
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
		SameSite: m.isSameSite(),
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
		SameSite: m.isSameSite(),
	}
}

func GoogleOAuthStateCookieName() string {
	return googleOAuthStateCookieName
}

func (m *CookieManager) isSecure() bool {
	return m.appCfg.Env != "development"
}

func (m *CookieManager) isSameSite() http.SameSite {
	if m.appCfg.Env == "production" || m.appCfg.Env == "prod" {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}
