package router

import (
	"shortivo/internal/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRoutes(
	g *gin.Engine,
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	linkHandler *handler.LinkHandler,
	analyticsHandler *handler.AnalyticsHandler,
	subscriptionHandler *handler.SubscriptionHandler,
) {
	r := g.Group("/api")

	authApi := r.Group("/auth")
	{
		authApi.POST("/signin", authHandler.SignIn)
		authApi.POST("/signup", authHandler.SignUp)
		authApi.POST("/refresh-token", authHandler.RefreshToken)
		authApi.POST("/logout", authHandler.SignOut)
	}

	// accountApi := r.Group("/accounts")
	// {

	// }

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
