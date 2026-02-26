package server

import (
	"context"
	"fmt"
	"net/http"
	"shortivo/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewHTTPServer(
	lc fx.Lifecycle,
	cfg config.AppConfig,
	ge *gin.Engine,
	log zerolog.Logger,
) *http.Server {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: ge,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info().
					Int("listen", cfg.Port).
					Msg("Starting listener")

				go func() {
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.Error().Err(err).Msg("http server failed")
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
