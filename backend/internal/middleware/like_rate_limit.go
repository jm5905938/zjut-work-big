package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
)

type LikeRateLimiter interface {
	Allow(
		ctx context.Context,
		userID uint64,
	) (bool, time.Duration, error)
}

func LikeRateLimit(limiter LikeRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get(ContextUserID)
		userID, valid := userIDValue.(uint64)
		if !exists || !valid {
			abortUnauthorized(c)
			return
		}

		allowed, retryAfter, err := limiter.Allow(
			c.Request.Context(),
			userID,
		)
		if err != nil {
			_ = c.Error(apperror.ServiceUnavailable(
				"点赞服务暂时不可用，请稍后重试",
				err,
			))
			c.Abort()
			return
		}

		if !allowed {
			retryAfterSeconds := int64(
				(retryAfter + time.Second - 1) / time.Second,
			)
			if retryAfterSeconds < 1 {
				retryAfterSeconds = 1
			}

			c.Header(
				"Retry-After",
				strconv.FormatInt(retryAfterSeconds, 10),
			)
			_ = c.Error(apperror.TooManyRequests(
				"操作过于频繁，请稍后再试",
			))
			c.Abort()
			return
		}

		c.Next()
	}
}
