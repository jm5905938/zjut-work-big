package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/middleware"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type PostHandler interface {
	Create(c *gin.Context)
}

type Dependencies struct {
	Auth   AuthHandler
	Posts  PostHandler
	Tokens middleware.TokenParser
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
	auth.POST("/login", deps.Auth.Login)

	posts := api.Group("/posts")
	posts.Use(middleware.Auth(deps.Tokens))
	posts.POST("", deps.Posts.Create)

	return r
}
