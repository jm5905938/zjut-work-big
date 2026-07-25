package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/model"
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
	_ = c.Error(apperror.Unauthorized("未登录或令牌无效"))
	c.Abort()
}

func RequireRole(requiredRole model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextUserRole)
		role, valid := roleValue.(model.UserRole)
		if !exists || !valid {
			abortUnauthorized(c)
			return
		}

		if role != requiredRole {
			_ = c.Error(apperror.Forbidden("当前用户无此权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}
