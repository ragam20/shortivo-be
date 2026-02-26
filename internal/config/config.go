package config

import "time"

type AppConfig struct {
	Env      string `mapstructure:"env"`
	Version  string `mapstructure:"version"`
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

type JwtConfig struct {
	AccessSecret string        `mapstructure:"access_secret"`
	Expiration   time.Duration `mapstructure:"expiration"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	TimeZone string `mapstructure:"time_zone"`
	SSLMode  string `mapstructure:"sslmode"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Redirect     string `mapstructure:"redirect"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JwtConfig
	GoogleOAuth GoogleOAuthConfig
}

func NewConfig(env *Env) *Config {
	return &Config{
		App:         env.App,
		Redis:       env.Redis,
		Database:    env.Database,
		GoogleOAuth: env.GoogleOAuth,
		JWT:         env.Jwt,
	}
}
