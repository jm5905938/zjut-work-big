package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			attributes := []any{
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", appErr.Status,
				"code", appErr.Code,
				"message", appErr.Message,
			}
			if appErr.Err != nil {
				attributes = append(attributes, "error", appErr.Err)
			}

			if appErr.Status >= http.StatusInternalServerError {
				logger.ErrorContext(
					c.Request.Context(),
					"application error",
					attributes...,
				)
			} else {
				logger.WarnContext(
					c.Request.Context(),
					"application error",
					attributes...,
				)
			}

			if c.Writer.Written() {
				return
			}

			response.Error(
				c,
				appErr.Status,
				appErr.Code,
				appErr.Message,
			)
			return
		}

		logger.ErrorContext(
			c.Request.Context(),
			"unhandled error",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"error", err,
		)

		if c.Writer.Written() {
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			500,
			"服务器死了啊）",
		)
	}
}
