package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort string

	MySQLHost     string
	MySQLPort     string
	MySQLDatabase string
	MySQLUser     string
	MySQLPassword string

	RedisHost      string
	RedisPort      string
	RedisPassword  string
	RedisDB        int
	LikeRateLimit  int64
	LikeRateWindow time.Duration

	JWTSecret    string
	JWTExpiresIn time.Duration

	AdminRegisterCode string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("MYSQL_HOST", "127.0.0.1")
	v.SetDefault("MYSQL_PORT", "3306")
	v.SetDefault("REDIS_HOST", "127.0.0.1")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("LIKE_RATE_LIMIT", 60)
	v.SetDefault("LIKE_RATE_WINDOW", "1m")
	v.SetDefault("JWT_EXPIRES_IN", "2h")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	cfg := &Config{
		AppPort: v.GetString("APP_PORT"),

		MySQLHost:     v.GetString("MYSQL_HOST"),
		MySQLPort:     v.GetString("MYSQL_PORT"),
		MySQLDatabase: v.GetString("MYSQL_DATABASE"),
		MySQLUser:     v.GetString("MYSQL_USER"),
		MySQLPassword: v.GetString("MYSQL_PASSWORD"),

		RedisHost:      v.GetString("REDIS_HOST"),
		RedisPort:      v.GetString("REDIS_PORT"),
		RedisPassword:  v.GetString("REDIS_PASSWORD"),
		RedisDB:        v.GetInt("REDIS_DB"),
		LikeRateLimit:  v.GetInt64("LIKE_RATE_LIMIT"),
		LikeRateWindow: v.GetDuration("LIKE_RATE_WINDOW"),

		JWTSecret:    v.GetString("JWT_SECRET"),
		JWTExpiresIn: v.GetDuration("JWT_EXPIRES_IN"),

		AdminRegisterCode: v.GetString("ADMIN_REGISTER_CODE"),
	}

	if cfg.MySQLDatabase == "" ||
		cfg.MySQLUser == "" ||
		cfg.MySQLPassword == "" {
		return nil, fmt.Errorf("MySQL 配置不完整")
	}

	if cfg.RedisPassword == "" {
		return nil, fmt.Errorf("Redis 配置不完整")
	}
	if cfg.LikeRateLimit <= 0 || cfg.LikeRateWindow <= 0 {
		return nil, fmt.Errorf("点赞限流配置无效")
	}

	if cfg.JWTSecret == "" || cfg.JWTExpiresIn <= 0 {
		return nil, fmt.Errorf("JWT 配置不完整")
	}

	return cfg, nil
}
