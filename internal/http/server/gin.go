package server

import (
	"time"

	"shortivo/internal/config"
	"shortivo/internal/http/limiter"
	"shortivo/internal/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinEngine(
	requestTimeout *middleware.RequestTimeoutMiddleware,
	errorHandling *middleware.ErrorMiddleware,
	requestLimtier *limiter.RequestLimiter,
	recovery *middleware.RecoveryMiddleware,
	appConfig *config.AppConfig,
) *gin.Engine {
	cors := cors.New(cors.Config{
		AllowOrigins:     isProduction(appConfig.Env),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Timezone"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r := gin.New()
	r.Use(
		gin.Logger(),
		recovery.Handle(),
		requestLimtier.HandleGlobal(),
		cors,
		requestTimeout.Handle(),
		errorHandling.Handle(),
	)

	return r
}

func isProduction(appEnv string) []string {
	if appEnv == "production" || appEnv == "prod" {
		return []string{"https://shortivo-fe.vercel.app"}
	}
	return []string{"http://localhost:5173"}
}
