package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindByID(ctx context.Context, id uint64) (*model.Post, error)
	List(
		ctx context.Context,
		page int,
		pageSize int,
	) ([]model.Post, int64, error)
	FindDetailByID(ctx context.Context, id uint64) (*model.Post, error)
	DeleteByIDAndAuthorID(
		ctx context.Context,
		postID uint64,
		authorID uint64,
	) (bool, error)
	ToggleLike(
		ctx context.Context,
		postID uint64,
		userID uint64,
	) (bool, error)
	FindLikedPostIDs(
		ctx context.Context,
		userID uint64,
		postIDs []uint64,
	) ([]uint64, error)
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

func (s *PostService) GetByID(
	ctx context.Context,
	postID uint64,
) (dto.PostDetailResponse, error) {
	post, err := s.posts.FindDetailByID(ctx, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PostDetailResponse{},
				apperror.NotFound("帖子不存在")
		}

		return dto.PostDetailResponse{}, apperror.Internal(err)
	}

	comments := make([]dto.CommentResponse, 0, len(post.Comments))
	for _, comment := range post.Comments {
		comments = append(comments, dto.CommentResponse{
			ID:      comment.ID,
			PostID:  comment.PostID,
			Content: comment.Content,
			Author: dto.UserResponse{
				ID:       comment.Author.ID,
				Username: comment.Author.Username,
				Name:     comment.Author.Name,
				Role:     comment.Author.Role,
			},
			CreatedAt: comment.CreatedAt,
		})
	}

	return dto.PostDetailResponse{
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
		Comments:     comments,
	}, nil
}

func (s *PostService) DeleteOwn(
	ctx context.Context,
	postID uint64,
	authorID uint64,
) error {
	post, err := s.posts.FindByID(ctx, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("帖子不存在")
		}

		return apperror.Internal(err)
	}

	if post.AuthorID != authorID {
		return apperror.Forbidden("无权删除他人的帖子")
	}

	deleted, err := s.posts.DeleteByIDAndAuthorID(
		ctx,
		postID,
		authorID,
	)
	if err != nil {
		return apperror.Internal(err)
	}
	if !deleted {
		return apperror.NotFound("帖子不存在")
	}

	return nil
}

func (s *PostService) ToggleLike(
	ctx context.Context,
	postID uint64,
	userID uint64,
) (dto.LikePostResponse, error) {
	isLiked, err := s.posts.ToggleLike(ctx, postID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.LikePostResponse{},
				apperror.NotFound("帖子不存在")
		}

		return dto.LikePostResponse{}, apperror.Internal(err)
	}

	return dto.LikePostResponse{
		PostID:  postID,
		IsLiked: isLiked,
	}, nil
}

func (s *PostService) GetLikeStatuses(
	ctx context.Context,
	userID uint64,
	req dto.GetLikeStatusesRequest,
) (dto.GetLikeStatusesResponse, error) {
	likedPostIDs, err := s.posts.FindLikedPostIDs(
		ctx,
		userID,
		req.PostIDs,
	)
	if err != nil {
		return dto.GetLikeStatusesResponse{}, apperror.Internal(err)
	}

	likedSet := make(map[uint64]struct{}, len(likedPostIDs))
	for _, postID := range likedPostIDs {
		likedSet[postID] = struct{}{}
	}

	statuses := make([]dto.PostLikeStatus, 0, len(req.PostIDs))
	for _, postID := range req.PostIDs {
		_, liked := likedSet[postID]
		statuses = append(statuses, dto.PostLikeStatus{
			PostID: postID,
			Liked:  liked,
		})
	}

	return dto.GetLikeStatusesResponse{
		Status: statuses,
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
