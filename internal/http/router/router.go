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
	analyticsHandler *handler.AnalyticsHandler,
) {
	api := g.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    http.StatusOK,
				"message":   "success",
				"timestamp": time.Now().UTC(),
				"data": gin.H{
					"message": "PONG",
				},
			})
		})

		api.GET("/debug/ip", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    http.StatusOK,
				"message":   "success",
				"timestamp": time.Now().UTC(),
				"data": gin.H{
					"clientIp": c.ClientIP(),
				},
			})
		})

		authAPI := api.Group("/auth")
		{
			authAPI.GET("/google", authHandler.GoogleSignIn)
			authAPI.GET("/google/redirect", authHandler.GoogleCallback)
			authAPI.POST("/signin", authHandler.SignIn)
			authAPI.POST("/signup", authHandler.SignUp)
			authAPI.POST("/refresh-token", authHandler.RefreshToken)
			authAPI.POST("/logout", accessTokenMiddleware.Handle(), authHandler.SignOut)
		}

		accountAPI := api.Group("/users", accessTokenMiddleware.Handle())
		{
			accountAPI.GET("/account", accountHandler.GetAccount)
		}

		linkAPI := api.Group("/links")
		{
			linkAPI.GET("", accessTokenMiddleware.Handle(), linkHandler.GetLinks)
			linkAPI.POST("", accessTokenMiddleware.Handle(), linkHandler.CreateShortLink)
			linkAPI.GET("/:linkId", accessTokenMiddleware.Handle(), linkHandler.GetLinkByID)
			linkAPI.PATCH("/:linkId", accessTokenMiddleware.Handle(), linkHandler.UpdateLink)
			linkAPI.DELETE("/:linkId", accessTokenMiddleware.Handle(), linkHandler.DeleteLink)
			linkAPI.GET("/monthly-limit", accessTokenMiddleware.Handle(), linkHandler.GetMonthlyLimit)
		}

		subscriptionAPI := api.Group("/subscriptions")
		{
			subscriptionAPI.GET("/current", accessTokenMiddleware.Handle(), subscriptionHandler.GetCurrentSubscription)
			subscriptionAPI.GET("/plan", accessTokenMiddleware.Handle(), subscriptionHandler.GetCurrentPlan)
			subscriptionAPI.POST("/cancel/:id", accessTokenMiddleware.Handle(), subscriptionHandler.CancelSubsription)
		}

		analyticsAPI := api.Group("/analytics", accessTokenMiddleware.Handle())
		{
			analyticsAPI.GET("/clicks", analyticsHandler.GetTotalClicks)
			analyticsAPI.GET("/link", analyticsHandler.GetTotalCreatedLink)
			analyticsAPI.GET("/:linkId/clicks", analyticsHandler.GetTotalClicksByLinkID)
			analyticsAPI.GET("/:linkId/referrer", analyticsHandler.GetEnagementByReferrer)
			analyticsAPI.GET("/:linkId/device", analyticsHandler.GetEngagementByDevices)
			analyticsAPI.GET("/:linkId/overtime", analyticsHandler.GetTotalEngagementByOvertime)
		}
	}

	g.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	g.GET("/:code", linkHandler.RedirectShortLink)

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
