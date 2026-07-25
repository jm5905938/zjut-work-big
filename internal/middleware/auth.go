package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/token"
)

const (
	ContextUserID   = "userID"
	ContextUsername = "username"
	ContextUserRole = "userRole"
)

type TokenParser interface {
	Parse(signedToken string) (*token.Claims, error)
}

func Auth(tokens TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		parts := strings.Fields(authorization)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			abortUnauthorized(c)
			return
		}

		claims, err := tokens.Parse(parts[1])
		if err != nil {
			abortUnauthorized(c)
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextUserRole, claims.Role)

		c.Next()
	}
}

func abortUnauthorized(c *gin.Context) {
	_ = c.Error(apperror.Unauthorized("登录状态无效或已过期"))
	c.Abort()
}
