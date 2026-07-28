package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/jm5905938/zjut-work-big/internal/cache"
	"github.com/jm5905938/zjut-work-big/internal/config"
	"github.com/jm5905938/zjut-work-big/internal/database"
	"github.com/jm5905938/zjut-work-big/internal/handler"
	"github.com/jm5905938/zjut-work-big/internal/ratelimit"
	"github.com/jm5905938/zjut-work-big/internal/repository"
	"github.com/jm5905938/zjut-work-big/internal/router"
	"github.com/jm5905938/zjut-work-big/internal/service"
	"github.com/jm5905938/zjut-work-big/internal/token"
)

func main() {
	binding.EnableDecoderDisallowUnknownFields = true

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenMySQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	redisClient, err := cache.OpenRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = redisClient.Close()
	}()

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)
	tokenManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresIn)
	authService := service.NewAuthService(userRepository, tokenManager)
	postService := service.NewPostService(postRepository)
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	likeLimiter := ratelimit.NewLikeLimiter(
		redisClient,
		cfg.LikeRateLimit,
		cfg.LikeRateWindow,
	)

	r := router.New(router.Dependencies{
		Auth:        authHandler,
		Posts:       postHandler,
		Tokens:      tokenManager,
		LikeLimiter: likeLimiter,
	})

	address := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("服务器启动：http://127.0.0.1%s", address)

	if err := r.Run(address); err != nil {
		log.Fatal("服务器启动失败：", err)
	}
}
