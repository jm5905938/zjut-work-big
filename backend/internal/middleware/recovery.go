package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		response.Error(
			c,
			http.StatusInternalServerError,
			500,
			"服务器死了啊）",
		)
	})
}
