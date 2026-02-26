package cache

import (
	"fmt"
	"shortivo/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisConfig(cfg *config.Config) config.RedisConfig {
	return cfg.Redis
}

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
}
