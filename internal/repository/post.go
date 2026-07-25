package repository

import (
	"context"

	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
)

const postWithCountsSelect = `posts.*,
	(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
	(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count`

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

func (r *PostRepository) List(
	ctx context.Context,
	page int,
	pageSize int,
) ([]model.Post, int64, error) {
	posts := make([]model.Post, 0)
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Post{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Select(postWithCountsSelect).
		Preload("Author").
		Order("posts.created_at DESC").
		Order("posts.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepository) FindDetailByID(
	ctx context.Context,
	id uint64,
) (*model.Post, error) {
	var post model.Post

	err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Select(postWithCountsSelect).
		Preload("Author").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("comments.created_at ASC").
				Order("comments.id ASC")
		}).
		Preload("Comments.Author").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}
