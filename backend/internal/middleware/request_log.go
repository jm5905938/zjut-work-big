package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()

		c.Next()

		logger.InfoContext(
			c.Request.Context(),
			"request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", float64(time.Since(startedAt).Microseconds())/1000,
			"response_size", c.Writer.Size(),
			"remote_addr", c.Request.RemoteAddr,
		)
	}
}
