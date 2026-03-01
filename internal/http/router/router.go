package router

import (
	"net/http"
	"shortivo/internal/http/handler"
	"shortivo/internal/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRoutes(
	g *gin.Engine,
	accessTokenMiddleware *middleware.AccessTokenMiddleware,
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	// linkHandler *handler.LinkHandler,
	// analyticsHandler *handler.AnalyticsHandler,
	// subscriptionHandler *handler.SubscriptionHandler,
) {
	r := g.Group("/api")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"message":   "pong",
			"timestamp": time.Now().UTC(),
		})
	})

	authApi := r.Group("/auth")
	{
		authApi.POST("/signin", authHandler.SignIn)
		authApi.POST("/signup", authHandler.SignUp)
		authApi.POST("/refresh-token", authHandler.RefreshToken)
		authApi.POST("/logout", accessTokenMiddleware.Handle(), authHandler.SignOut)
	}

	accountApi := r.Group("/accounts", accessTokenMiddleware.Handle())
	{
		accountApi.GET("/profile", accountHandler.GetProfile)
	}

	// linkApi := r.Group("/links")
	// {

	// }

	// analyticsApi := r.Group("/analytics")
	// {

	// }

	// subscriptionApi := r.Group("/subscriptions")
	// {

	// }
}
