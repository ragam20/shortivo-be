package server

import (
	"net/http"
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
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":    http.StatusNotFound,
			"message":   "route not found",
			"timestamp": time.Now().UTC(),
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":    http.StatusMethodNotAllowed,
			"message":   "method no allowed",
			"timestamp": time.Now().UTC(),
		})
	})

	return r
}
