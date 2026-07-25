package repository

import (
	"context"

	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
