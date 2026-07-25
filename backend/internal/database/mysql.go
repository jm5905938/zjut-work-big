package database

import (
	"context"
	"fmt"
	"net"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jm5905938/zjut-work-big/internal/config"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(cfg *config.Config) (*gorm.DB, error) {
	mysqlConfig := mysqldriver.Config{
		User:      cfg.MySQLUser,
		Passwd:    cfg.MySQLPassword,
		Net:       "tcp",
		Addr:      net.JoinHostPort(cfg.MySQLHost, cfg.MySQLPort),
		DBName:    cfg.MySQLDatabase,
		ParseTime: true,
		Loc:       time.Local,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	dsn := mysqlConfig.FormatDSN()

	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("打开 MySQL 失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	return db, nil
}
