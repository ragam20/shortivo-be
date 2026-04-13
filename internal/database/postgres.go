package database

import (
	"shortivo/internal/config"
	"shortivo/internal/database/entity"

	appLogger "shortivo/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(cfg *config.DatabaseConfig, appCfg *config.AppConfig, log appLogger.Logger) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true,
	}

	pgConfig := postgres.Config{
		DSN: cfg.URL,
	}

	if appCfg.Env == "production" {
		gormCfg.PrepareStmt = false
		pgConfig.PreferSimpleProtocol = true
	}

	db, err := gorm.Open(postgres.New(pgConfig), gormCfg)
	if err != nil {
		log.Error("failed to connect database", err, map[string]any{})
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
			log.Error("failed to migrate", err, map[string]any{})
			return nil, err
		}
	}

	return db, nil
}
