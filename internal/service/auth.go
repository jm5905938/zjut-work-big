package service

import (
	"context"
	"errors"

	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/model"
	"github.com/jm5905938/zjut-work-big/internal/password"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)
}

type AuthService struct {
	users UserRepository
}

func NewAuthService(users UserRepository) *AuthService {
	return &AuthService{
		users: users,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (dto.UserResponse, error) {
	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return dto.UserResponse{}, apperror.Internal(err)
	}

	user := &model.User{
		Username:     req.Username,
		Name:         req.Name,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return dto.UserResponse{},
				apperror.Conflict("用户名已存在")
		}

		return dto.UserResponse{}, apperror.Internal(err)
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}, nil
}
