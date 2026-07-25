package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/middleware"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

type PostService interface {
	Create(
		ctx context.Context,
		authorID uint64,
		req dto.CreatePostRequest,
	) (dto.PostResponse, error)
	List(
		ctx context.Context,
		query dto.ListPostsQuery,
	) (dto.PostListResponse, error)
}

type PostHandler struct {
	posts PostService
}

func NewPostHandler(posts PostService) *PostHandler {
	return &PostHandler{
		posts: posts,
	}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("参数校验失败"))
		return
	}

	authorIDValue, exists := c.Get(middleware.ContextUserID)
	authorID, valid := authorIDValue.(uint64)
	if !exists || !valid {
		_ = c.Error(apperror.Unauthorized("未登录或令牌无效"))
		return
	}

	post, err := h.posts.Create(
		c.Request.Context(),
		authorID,
		req,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, post)
}

func (h *PostHandler) List(c *gin.Context) {
	query := dto.ListPostsQuery{
		Page:     1,
		PageSize: 20,
		Sort:     "latest",
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(apperror.BadRequest("参数校验失败"))
		return
	}

	posts, err := h.posts.List(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, posts)
}
