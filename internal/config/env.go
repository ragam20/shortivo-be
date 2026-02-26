package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Env struct {
	App         AppConfig         `mapstructure:"app"`
	Jwt         JwtConfig         `mapstructure:"jwt"`
	Database    DatabaseConfig    `mapstructure:"database"`
	GoogleOAuth GoogleOAuthConfig `mapstructure:"google_oauth"`
	Redis       RedisConfig       `mapstructure:"redis"`
}

func NewEnv() (*Env, error) {
	configName := os.Getenv("APP_ENV")
	if configName == "" {
		configName = "development"
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("internal/config")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %q: %w", configName, err)
	}

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", configName, err)
	}

	return &env, nil
}
