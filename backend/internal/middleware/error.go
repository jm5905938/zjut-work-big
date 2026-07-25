package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			response.Error(
				c,
				appErr.Status,
				appErr.Code,
				appErr.Message,
			)
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
