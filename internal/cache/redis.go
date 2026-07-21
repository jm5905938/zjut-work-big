package cache

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jm5905938/zjut-work-big/internal/config"
	"github.com/redis/go-redis/v9"
)

func OpenRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(cfg.RedisHost, cfg.RedisPort),

		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return client, nil
}
