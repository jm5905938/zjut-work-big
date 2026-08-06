package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/middleware"
	"github.com/jm5905938/zjut-work-big/internal/model"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type PostHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	DeleteOwn(c *gin.Context)
	DeleteAny(c *gin.Context)
	ToggleLike(c *gin.Context)
	GetLikeStatuses(c *gin.Context)
	CreateComment(c *gin.Context)
}

type Dependencies struct {
	Auth        AuthHandler
	Posts       PostHandler
	Tokens      middleware.TokenParser
	LikeLimiter middleware.LikeRateLimiter
	Logger      *slog.Logger
}

func New(deps Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.RequestLog(deps.Logger),
		middleware.Recovery(deps.Logger),
		middleware.ErrorHandler(deps.Logger),
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
	posts.GET("", deps.Posts.List)
	posts.GET("/:post_id", deps.Posts.GetByID)
	posts.DELETE("/:post_id", deps.Posts.DeleteOwn)
	posts.POST(
		"/:post_id/like",
		middleware.LikeRateLimit(deps.LikeLimiter),
		deps.Posts.ToggleLike,
	)
	posts.POST("/likes", deps.Posts.GetLikeStatuses)
	posts.POST("/:post_id/comment", deps.Posts.CreateComment)

	admin := api.Group("/admin")
	admin.Use(
		middleware.Auth(deps.Tokens),
		middleware.RequireRole(model.UserRoleAdmin),
	)
	admin.DELETE("/posts/:post_id", deps.Posts.DeleteAny)

	return r
}
