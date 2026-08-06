package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(
					c.Request.Context(),
					"panic recovered",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"panic", fmt.Sprint(recovered),
					"stack", string(debug.Stack()),
				)

				response.Error(
					c,
					http.StatusInternalServerError,
					500,
					"服务器死了啊）",
				)
				c.Abort()
			}
		}()

		c.Next()
	}
}
