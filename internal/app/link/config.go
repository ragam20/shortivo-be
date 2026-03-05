package link

import "shortivo/internal/config"

type LinkConfig struct {
	ShortenerBaseURL string
}

func NewLinkConfig(cfg *config.AppConfig) *LinkConfig {
	return &LinkConfig{ShortenerBaseURL: cfg.ShortenerBaseURL}
}
