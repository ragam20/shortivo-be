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
	"shortivo/internal/http/limiter"
	"shortivo/internal/http/middleware"
	"shortivo/internal/http/router"
	"shortivo/internal/http/server"
	"shortivo/internal/idgen"
	"shortivo/internal/logger"

	"go.uber.org/fx"
)

func main() {
	configModule := fx.Module("config",
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

	databaseModule := fx.Module("database",
		fx.Provide(
			database.NewPostgresConnection,
		),
	)

	infraModule := fx.Module("infra",
		fx.Provide(
			idgen.NewSnowflakeGen,
			logger.NewLogger,
			cache.NewRedisClient,
			cache.NewStore,
		),
	)

	appHTTPModule := fx.Module("apphttp",
		fx.Provide(
			limiter.NewRequestLimiter,
			middleware.NewRequestTimeoutMiddleware,
			middleware.NewErrorMiddleware,
			middleware.NewRecoveryMiddleware,

			server.NewGinEngine,
			server.NewHTTPServer,

			appHttp.NewCookieManager,
		),
	)

	appRepoModule := fx.Module("appRepo",
		fx.Provide(
			repo.NewAuthRepository,
			repo.NewAccountRepository,
			repo.NewLinkRepository,
			repo.NewPlanRespository,
			repo.NewAnalyticsRepository,
			repo.NewSubscriptionRepository,
		),
	)

	authModule := fx.Module("auth",
		fx.Provide(
			auth.NewJWTManager,

			middleware.NewAccessTokenMiddleware,

			oauth.NewGoogleOAuth,
			oauth.NewGoogleOAuthProvider,

			auth.NewService,
			handler.NewAuthHandler,
		),
	)

	accountModule := fx.Module("account",
		fx.Provide(
			account.NewService,
			handler.NewAccountHandler,
		),
	)

	linkModule := fx.Module("link",
		fx.Provide(
			link.NewLinkConfig,
			link.NewLinkCache,
			link.NewLinkService,
			handler.NewLinkHandler,
		),
	)

	analyticsModule := fx.Module("analytics",
		fx.Provide(
			analytics.NewService,
			handler.NewAnalyticsHandler,
		),
	)

	subscriptionModule := fx.Module("subscription",
		fx.Provide(
			subscription.NewService,
			handler.NewSubscriptionHandler,
		),
	)

	app := fx.New(fx.Options(
		configModule,
		databaseModule,
		infraModule,
		appHTTPModule,
		appRepoModule,
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
