package database

import (
	"fmt"
	"shortivo/internal/config"
	"shortivo/internal/database/entity"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(cfg *config.DatabaseConfig, appCfg *config.AppConfig, log *zerolog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect database")
		return nil, err
	}

	if appCfg.Env == "development" {
		if err := db.AutoMigrate(
			entity.UserEntity{},
			entity.RefreshTokenEntity{},
			entity.PlanEntity{},
			entity.LinkEntity{},
			entity.LinkClicksEntity{},
			entity.SubscriptionEntity{},
		); err != nil {
			return nil, err
		}
	}

	return db, nil
}
