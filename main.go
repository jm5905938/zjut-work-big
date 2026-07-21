package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/jm5905938/zjut-work-big/internal/cache"
	"github.com/jm5905938/zjut-work-big/internal/config"
	"github.com/jm5905938/zjut-work-big/internal/database"
	"github.com/jm5905938/zjut-work-big/internal/router"
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

	r := router.New(router.Dependencies{
		DB:    db,
		Redis: redisClient,
	})

	address := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("服务器启动：http://127.0.0.1%s", address)

	if err := r.Run(address); err != nil {
		log.Fatal("服务器启动失败：", err)
	}
}
