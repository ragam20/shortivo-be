package logger

import (
	"fmt"
	"os"
	"shortivo/internal/config"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.AppConfig) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("invalid log level: %w", err)
	}

	logger := zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()

	if cfg.Env == "development" {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out: os.Stdout,
		})
	}

	return logger, nil
}
