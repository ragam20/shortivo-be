package router

import (
	"net/http"
	"time"

	"shortivo/internal/http/handler"
	"shortivo/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRoutes(
	g *gin.Engine,
	accessTokenMiddleware *middleware.AccessTokenMiddleware,
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	linkHandler *handler.LinkHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	// analyticsHandler *handler.AnalyticsHandler,
) {
	g.GET("/r/:code", linkHandler.RedirectShortLink)
	g.HEAD("/r/:code", linkHandler.RedirectShortLink)

	api := g.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    http.StatusOK,
				"message":   "pong",
				"timestamp": time.Now().UTC(),
			})
		})

		authApi := api.Group("/auth")
		{
			authApi.POST("/signin", authHandler.SignIn)
			authApi.POST("/signup", authHandler.SignUp)
			authApi.POST("/refresh-token", authHandler.RefreshToken)
			authApi.POST("/logout", accessTokenMiddleware.Handle(), authHandler.SignOut)
		}

		accountApi := api.Group("/accounts", accessTokenMiddleware.Handle())
		{
			accountApi.GET("/profile", accountHandler.GetProfile)
		}

		linkApi := api.Group("/links")
		{
			linkApi.GET("", accessTokenMiddleware.Handle(), linkHandler.GetLinks)
			linkApi.POST("/create", accessTokenMiddleware.Handle(), linkHandler.CreateShortLink)
			linkApi.PATCH("/:linkId", accessTokenMiddleware.Handle(), linkHandler.UpdateLink)
			linkApi.DELETE("/:linkId", accessTokenMiddleware.Handle(), linkHandler.DeleteLink)
		}

		subscriptionApi := api.Group("/subscriptions")
		{
			subscriptionApi.POST("/cancel/:id", accessTokenMiddleware.Handle(), subscriptionHandler.CancelSubsription)
		}
	}

	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":    http.StatusNotFound,
			"message":   "route not found",
			"timestamp": time.Now().UTC(),
		})
	})

	g.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":    http.StatusMethodNotAllowed,
			"message":   "method no allowed",
			"timestamp": time.Now().UTC(),
		})
	})
}
