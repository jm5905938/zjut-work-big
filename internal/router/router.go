package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/middleware"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

type AuthHandler interface {
	Register(c *gin.Context)
}

type Dependencies struct {
	Auth AuthHandler
}

func New(deps Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(
		gin.Logger(),
		middleware.Recovery(),
		middleware.ErrorHandler(),
	)

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"status": "healthy",
			"mysql":  "connected",
			"redis":  "connected",
		})
	})

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", deps.Auth.Register)

	return r
}
