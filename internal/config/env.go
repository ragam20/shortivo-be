package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Env struct {
	App         AppConfig         `mapstructure:"app"`
	Jwt         JwtConfig         `mapstructure:"jwt"`
	Database    DatabaseConfig    `mapstructure:"db"`
	GoogleOAuth GoogleOAuthConfig `mapstructure:"google_oauth"`
	Redis       RedisConfig       `mapstructure:"redis"`
}

func NewEnv() (*Env, error) {
	appEnv := os.Getenv("APP_ENV")

	configName := ""
	if appEnv == "" {
		configName = "development"
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("internal/config")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if !isProduction(appEnv) {
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if !errors.As(err, &notFound) {
				return nil, fmt.Errorf("read config %q: %w", configName, err)
			}
		}
	}

	if err := bindEnvs(v, getKeys()); err != nil {
		return nil, fmt.Errorf("bind env: %w", err)
	}

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", configName, err)
	}
	if err := env.Validate(); err != nil {
		return nil, fmt.Errorf("validate config %q: %w", configName, err)
	}

	data, _ := json.Marshal(map[string]any{
		"level": "info",
		"env":   env.App.Env,
		"port":  env.App.Port,
		"msg":   "config load",
	})
	fmt.Println(string(data))

	return &env, nil
}

func isProduction(appEnv string) bool {
	return appEnv == "production" || appEnv == "prod"
}

func bindEnvs(v *viper.Viper, keys []string) error {
	for _, k := range keys {
		if err := v.BindEnv(k); err != nil {
			return err
		}
	}
	return nil
}

func getKeys() []string {
	return []string{
		"app.env",
		"app.port",
		"app.version",
		"app.log_level",
		"app.request_timeout",

		"database.url",

		"jwt.access_secret",
		"jwt.expiration",
		"jwt.refresh_exp",

		"redis.url",
	}
}
