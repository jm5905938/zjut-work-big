package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin/binding"
	"github.com/jm5905938/zjut-work-big/internal/cache"
	"github.com/jm5905938/zjut-work-big/internal/config"
	"github.com/jm5905938/zjut-work-big/internal/database"
	"github.com/jm5905938/zjut-work-big/internal/handler"
	"github.com/jm5905938/zjut-work-big/internal/logger"
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

	appLogger, logFile, err := logger.New(cfg.LogFile, cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logFile.Close()
	}()
	slog.SetDefault(appLogger)

	db, err := database.OpenMySQL(cfg)
	if err != nil {
		appLogger.Error("MySQL 初始化失败", "error", err)
		return
	}

	if err := database.Migrate(db); err != nil {
		appLogger.Error("数据库迁移失败", "error", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("获取数据库连接失败", "error", err)
		return
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	redisClient, err := cache.OpenRedis(cfg)
	if err != nil {
		appLogger.Error("Redis 初始化失败", "error", err)
		return
	}
	defer func() {
		_ = redisClient.Close()
	}()

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)
	tokenManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresIn)
	authService := service.NewAuthService(
		userRepository,
		tokenManager,
		cfg.AdminRegisterCode,
	)
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
		Logger:      appLogger,
	})

	address := fmt.Sprintf(":%s", cfg.AppPort)
	appLogger.Info(
		"服务器启动",
		"address", fmt.Sprintf("http://127.0.0.1%s", address),
	)

	if err := r.Run(address); err != nil {
		appLogger.Error("服务器启动失败", "error", err)
	}
}
