package limiter

import (
	"net/http"

	"shortivo/internal/http/response"
	"shortivo/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	liblimiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RequestLimiter struct {
	rdsClient *redis.Client
	log       logger.Logger
}

func NewRequestLimiter(rdsClient *redis.Client, log logger.Logger) *RequestLimiter {
	return &RequestLimiter{
		rdsClient: rdsClient,
		log:       log,
	}
}

func (r *RequestLimiter) HandleGlobal() gin.HandlerFunc {
	rate, err := liblimiter.NewRateFromFormatted("9-s")
	if err != nil {
		r.log.Error("request limiter err", err, nil)
		panic(err)
	}

	store, err := sredis.NewStoreWithOptions(r.rdsClient, liblimiter.StoreOptions{
		Prefix:   "limit_global_api",
		MaxRetry: 3,
	})
	if err != nil {
		r.log.Error("request limiter err", err, nil)
		panic(err)
	}

	return mgin.NewMiddleware(
		liblimiter.New(store, rate),
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.TooManyRequestError())
		}),
		mgin.WithErrorHandler(func(c *gin.Context, err error) {
			r.log.Error("request limiter err", err, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		}),
	)
}
