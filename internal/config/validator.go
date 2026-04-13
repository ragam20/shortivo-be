package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
)

var validDBSSLModes = []string{
	"disable",
	"allow",
	"prefer",
	"require",
	"verify-ca",
	"verify-full",
}

func (e *Env) Validate() error {
	var errs []string

	validateAppConfig(e.App, &errs)
	validateDatabaseConfig(e.Database, &errs)
	validateJWTConfig(e.Jwt, &errs)
	validateRedisConfig(e.Redis, &errs)
	validateGoogleOAuthConfig(e.GoogleOAuth, &errs)

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("invalid config: %s", strings.Join(errs, "; "))
}

func validateAppConfig(cfg AppConfig, errs *[]string) {
	if strings.TrimSpace(cfg.Env) == "" {
		*errs = append(*errs, "app.env is required")
	}

	if strings.TrimSpace(cfg.Version) == "" {
		*errs = append(*errs, "app.version is required")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		*errs = append(*errs, "app.port must be between 1 and 65535")
	}

	if strings.TrimSpace(cfg.LogLevel) == "" {
		*errs = append(*errs, "app.log_level is required")
	} else if _, err := zerolog.ParseLevel(cfg.LogLevel); err != nil {
		*errs = append(*errs, "app.log_level is invalid")
	}

	if cfg.NodeID < 0 {
		*errs = append(*errs, "app.node_id must be >= 0")
	}

	if cfg.RequestTimeout <= 0 {
		*errs = append(*errs, "app.request_timeout must be > 0")
	}
}

func validateDatabaseConfig(cfg DatabaseConfig, errs *[]string) {
	if strings.TrimSpace(cfg.URL) == "" {
		*errs = append(*errs, "database.url is required")
	}
}

func validateJWTConfig(cfg JwtConfig, errs *[]string) {
	if strings.TrimSpace(cfg.AccessSecret) == "" {
		*errs = append(*errs, "jwt.access_secret is required")
	}

	if cfg.Expiration <= 0 {
		*errs = append(*errs, "jwt.expiration must be > 0")
	}

	if cfg.RefreshExp <= 0 {
		*errs = append(*errs, "jwt.refresh_exp must be > 0")
	}
}

func validateRedisConfig(cfg RedisConfig, errs *[]string) {
	if strings.TrimSpace(cfg.URL) == "" {
		*errs = append(*errs, "redis.url is required")
	}
}

func validateGoogleOAuthConfig(cfg GoogleOAuthConfig, errs *[]string) {
	usesOAuth := strings.TrimSpace(cfg.ClientID) != "" ||
		strings.TrimSpace(cfg.ClientSecret) != "" ||
		strings.TrimSpace(cfg.Redirect) != ""

	if !usesOAuth {
		return
	}

	if strings.TrimSpace(cfg.ClientID) == "" {
		*errs = append(*errs, "google_oauth.client_id is required when google oauth is enabled")
	}

	if strings.TrimSpace(cfg.ClientSecret) == "" {
		*errs = append(*errs, "google_oauth.client_secret is required when google oauth is enabled")
	}

	if strings.TrimSpace(cfg.Redirect) == "" {
		*errs = append(*errs, "google_oauth.redirect is required when google oauth is enabled")
		return
	}

	redirectURL, err := url.ParseRequestURI(cfg.Redirect)
	if err != nil || redirectURL.Scheme == "" || redirectURL.Host == "" {
		*errs = append(*errs, "google_oauth.redirect must be a valid absolute URL")
	}
}
