package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/middleware"
	"github.com/jm5905938/zjut-work-big/internal/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Redis *redis.Client
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

	return r
}
