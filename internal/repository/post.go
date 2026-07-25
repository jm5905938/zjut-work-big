package repository

import (
	"context"

	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *PostRepository) DeleteByIDAndAuthorID(
	ctx context.Context,
	postID uint64,
	authorID uint64,
) (bool, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", postID, authorID).
		Delete(&model.Post{})
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *PostRepository) ToggleLike(
	ctx context.Context,
	postID uint64,
	userID uint64,
) (bool, error) {
	var isLiked bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post model.Post
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").
			First(&post, postID).Error; err != nil {
			return err
		}

		var like model.PostLike
		result := tx.
			Where("post_id = ? AND user_id = ?", postID, userID).
			Limit(1).
			Find(&like)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			like = model.PostLike{
				PostID: postID,
				UserID: userID,
			}
			if err := tx.Create(&like).Error; err != nil {
				return err
			}

			isLiked = true
			return nil
		}

		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		isLiked = false
		return nil
	})
	if err != nil {
		return false, err
	}

	return isLiked, nil
}

func (r *PostRepository) FindLikedPostIDs(
	ctx context.Context,
	userID uint64,
	postIDs []uint64,
) ([]uint64, error) {
	likedPostIDs := make([]uint64, 0)

	err := r.db.WithContext(ctx).
		Model(&model.PostLike{}).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Pluck("post_id", &likedPostIDs).Error
	if err != nil {
		return nil, err
	}

	return likedPostIDs, nil
}

func (r *PostRepository) CreateComment(
	ctx context.Context,
	comment *model.Comment,
) (*model.Comment, error) {
	var createdComment model.Comment

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post model.Post
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").
			First(&post, comment.PostID).Error; err != nil {
			return err
		}

		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		return tx.
			Preload("Author").
			First(&createdComment, comment.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return &createdComment, nil
}
