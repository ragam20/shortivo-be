package server

import "shortivo/internal/config"

func NewAppConfig(cfg *config.Config) config.AppConfig {
	return cfg.App
}
