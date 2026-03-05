package server

import (
	"shortivo/internal/http/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinEngine(
	requestTimeout *middleware.RequestTimeoutMiddleware,
	errorHandling *middleware.ErrorMiddleware,
	recovery *middleware.RecoveryMiddleware,
) *gin.Engine {
	cors := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r := gin.New()
	r.Use(
		gin.Logger(),
		recovery.Handle(),
		cors,
		requestTimeout.Handle(),
		errorHandling.Handle(),
	)

	return r
}
