package database

import (
	"fmt"

	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
	); err != nil {
		return fmt.Errorf("迁移数据库失败: %w", err)
	}

	return nil
}
