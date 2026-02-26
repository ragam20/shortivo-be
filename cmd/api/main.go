package api

import (
	"shortivo/internal/app/account"
	"shortivo/internal/app/analytics"
	"shortivo/internal/app/auth"
	"shortivo/internal/app/auth/oauth"
	"shortivo/internal/app/link"
	"shortivo/internal/app/subscription"
	"shortivo/internal/cache"
	"shortivo/internal/config"
	"shortivo/internal/database"
	dbRepo "shortivo/internal/database/repo"
	"shortivo/internal/http/server"
	"shortivo/internal/logger"

	"go.uber.org/fx"
)

func main() {
	var databaseModule = fx.Module("database",
		fx.Provide(
			database.NewDatabaseConfig,
			database.NewPostgresConnection,
			dbRepo.NewAuthRepository,
			dbRepo.NewAccountRepository,
			dbRepo.NewLinkRepository,
			dbRepo.NewSubscriptionRepository,
			dbRepo.NewAnalyticsRepository,
		),
	)

	var authModule = fx.Module("auth",
		fx.Provide(
			auth.NewJWTConfig,
			auth.NewJWTManager,
			oauth.NewGoogleOAuthConfig,
			oauth.NewGoogleOAuth,
			oauth.NewGoogleOAuthProvider,
			auth.NewService,
		),
	)

	var accountModule = fx.Module("account",
		fx.Provide(
			account.NewService,
		),
	)

	var linkModule = fx.Module("link",
		fx.Provide(
			link.NewService,
		),
	)

	var analyticsModule = fx.Module("analytics",
		fx.Provide(analytics.NewService),
	)

	var subscriptionModule = fx.Module("subscription",
		fx.Provide(subscription.NewService),
	)

	fx.New(fx.Options(
		fx.Provide(
			config.NewEnv,
			config.NewConfig,
			server.NewAppConfig,
			logger.NewLogger,
			cache.NewRedisConfig,
			cache.NewRedisClient,
			cache.NewStore,
		),

		databaseModule,
		authModule,
		accountModule,
		linkModule,
		analyticsModule,
		subscriptionModule,
	))
}
