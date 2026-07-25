package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/response"
)

type AuthService interface {
	Register(
		ctx context.Context,
		req dto.RegisterRequest,
	) (dto.UserResponse, error)
}

type AuthHandler struct {
	auth AuthService
}

func NewAuthHandler(auth AuthService) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(
			apperror.BadRequest("参数校验失败"),
		)
		return
	}

	user, err := h.auth.Register(
		c.Request.Context(),
		req,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, user)
}
