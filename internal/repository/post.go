package repository

import (
	"context"

	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(
	ctx context.Context,
	post *model.Post,
) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Post, error) {
	var post model.Post

	if err := r.db.WithContext(ctx).
		Preload("Author").
		First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
