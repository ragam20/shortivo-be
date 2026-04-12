package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"shortivo/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleOAuth(cfg *config.GoogleOAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.Redirect,

		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	EmailVerified bool   `json:"email_verified"`
}

type GoogleOAuthProvider interface {
	GetUser(ctx context.Context, code string) (*GoogleUser, error)
	GetURL(state string) string
	Enabled() bool
}

type googleOAuthProvider struct {
	cfg *oauth2.Config
}

func NewGoogleOAuthProvider(cfg *oauth2.Config) GoogleOAuthProvider {
	return googleOAuthProvider{cfg: cfg}
}

func (g googleOAuthProvider) GetUser(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("error get user from google oauth")
	}

	return g.getClient(ctx, token)
}

func (g googleOAuthProvider) GetURL(state string) string {
	return g.cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
}

func (g googleOAuthProvider) Enabled() bool {
	return g.cfg != nil &&
		g.cfg.ClientID != "" &&
		g.cfg.ClientSecret != "" &&
		g.cfg.RedirectURL != ""
}

func (g googleOAuthProvider) getClient(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := g.cfg.Client(ctx, token)

	res, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, errors.New("error get google client")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo request failed with status %d", res.StatusCode)
	}

	var gUser GoogleUser
	if err := json.NewDecoder(res.Body).Decode(&gUser); err != nil {
		return nil, errors.New("invalid decode google user")
	}

	return &gUser, nil
}
