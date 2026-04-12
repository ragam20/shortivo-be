package server

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"shortivo/internal/config"
	"shortivo/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewHTTPServer(
	lc fx.Lifecycle,
	cfg *config.AppConfig,
	ge *gin.Engine,
	log logger.Logger,
) *http.Server {
	port := ""
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	if port == "" {
		if cfg.Port != 0 {
			port = strconv.Itoa(cfg.Port)
		}
	}

	srv := &http.Server{
		Addr:    port,
		Handler: ge,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info(
					"Starting listener",
					map[string]any{
						"listen": cfg.Port,
					},
				)

				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.Error("http server failed", err, map[string]any{})
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		},
	)

	return srv
}
