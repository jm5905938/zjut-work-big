package service

import (
	"context"
	"strings"

	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/model"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindByID(ctx context.Context, id uint64) (*model.Post, error)
	List(
		ctx context.Context,
		page int,
		pageSize int,
	) ([]model.Post, int64, error)
}

type PostService struct {
	posts PostRepository
}

func NewPostService(posts PostRepository) *PostService {
	return &PostService{
		posts: posts,
	}
}

func (s *PostService) Create(
	ctx context.Context,
	authorID uint64,
	req dto.CreatePostRequest,
) (dto.PostResponse, error) {
	if strings.TrimSpace(req.Content) == "" {
		return dto.PostResponse{}, apperror.BadRequest("帖子内容不能为空")
	}

	post := &model.Post{
		AuthorID: authorID,
		Content:  req.Content,
	}

	if err := s.posts.Create(ctx, post); err != nil {
		return dto.PostResponse{}, apperror.Internal(err)
	}

	createdPost, err := s.posts.FindByID(ctx, post.ID)
	if err != nil {
		return dto.PostResponse{}, apperror.Internal(err)
	}

	return dto.PostResponse{
		ID:      createdPost.ID,
		Content: createdPost.Content,
		Author: dto.UserResponse{
			ID:       createdPost.Author.ID,
			Username: createdPost.Author.Username,
			Name:     createdPost.Author.Name,
			Role:     createdPost.Author.Role,
		},
		CreatedAt: createdPost.CreatedAt,
	}, nil
}

func (s *PostService) List(
	ctx context.Context,
	query dto.ListPostsQuery,
) (dto.PostListResponse, error) {
	posts, total, err := s.posts.List(
		ctx,
		query.Page,
		query.PageSize,
	)
	if err != nil {
		return dto.PostListResponse{}, apperror.Internal(err)
	}

	items := make([]dto.PostListItemResponse, 0, len(posts))
	for _, post := range posts {
		items = append(items, dto.PostListItemResponse{
			ID:      post.ID,
			Content: post.Content,
			Author: dto.UserResponse{
				ID:       post.Author.ID,
				Username: post.Author.Username,
				Name:     post.Author.Name,
				Role:     post.Author.Role,
			},
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    post.CreatedAt,
		})
	}

	return dto.PostListResponse{
		Items: items,
		Meta: dto.PaginationMeta{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    total,
		},
	}, nil
}
