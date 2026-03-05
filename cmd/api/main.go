package main

import (
	"net/http"
	"shortivo/internal/app/account"
	"shortivo/internal/app/analytics"
	"shortivo/internal/app/auth"
	"shortivo/internal/app/auth/oauth"
	"shortivo/internal/app/link"
	"shortivo/internal/app/subscription"
	"shortivo/internal/cache"
	"shortivo/internal/config"
	"shortivo/internal/database"
	"shortivo/internal/database/repo"
	appHttp "shortivo/internal/http"
	"shortivo/internal/http/handler"
	"shortivo/internal/http/middleware"
	"shortivo/internal/http/router"
	"shortivo/internal/http/server"
	"shortivo/internal/idgen"
	"shortivo/internal/logger"

	"go.uber.org/fx"
)

func main() {
	var configModule = fx.Module("config",
		fx.Provide(
			config.NewEnv,
			config.NewConfig,
			config.NewAppConfig,
			config.NewDatabaseConfig,
			config.NewRedisConfig,
			config.NewJWTConfig,
			config.NewGoogleOAuthConfig,
		),
	)

	var databaseModule = fx.Module("database",
		fx.Provide(
			database.NewPostgresConnection,
		),
	)

	var infraModule = fx.Module("infra",
		fx.Provide(
			idgen.NewSnowflakeGen,
			logger.NewLogger,
			cache.NewRedisClient,
			cache.NewStore,
		),
	)

	var appHttpModule = fx.Module("apphttp",
		fx.Provide(
			middleware.NewRequestTimeoutMiddleware,
			middleware.NewErrorMiddleware,
			middleware.NewRecoveryMiddleware,

			server.NewGinEngine,
			server.NewHTTPServer,

			appHttp.NewCookieManager,
		),
	)

	var authModule = fx.Module("auth",
		fx.Provide(
			auth.NewJWTManager,

			middleware.NewAccessTokenMiddleware,

			oauth.NewGoogleOAuth,
			oauth.NewGoogleOAuthProvider,

			repo.NewAuthRepository,
			auth.NewService,
			handler.NewAuthHandler,
		),
	)

	var accountModule = fx.Module("account",
		fx.Provide(
			repo.NewAccountRepository,
			account.NewService,
			handler.NewAccountHandler,
		),
	)

	var linkModule = fx.Module("link",
		fx.Provide(
			link.NewLinkConfig,
			link.NewLinkCache,
			repo.NewLinkRepository,
			link.NewLinkService,
			handler.NewLinkHandler,
		),
	)

	var analyticsModule = fx.Module("analytics",
		fx.Provide(
			repo.NewAnalyticsRepository,
			analytics.NewService,
			handler.NewAnalyticsHandler,
		),
	)

	var subscriptionModule = fx.Module("subscription",
		fx.Provide(
			repo.NewSubscriptionRepository,
			subscription.NewService,
			handler.NewSubscriptionHandler,
		),
	)

	app := fx.New(fx.Options(
		configModule,
		databaseModule,
		infraModule,
		appHttpModule,
		authModule,
		accountModule,
		linkModule,
		analyticsModule,
		subscriptionModule,

		fx.Invoke(
			func(_ *http.Server) {},
			router.NewRoutes,
		),
	))

	app.Run()
}
