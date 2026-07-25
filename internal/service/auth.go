package service

import (
	"context"
	"errors"

	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/model"
	"github.com/jm5905938/zjut-work-big/internal/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)
}

type TokenGenerator interface {
	Generate(user *model.User) (string, int64, error)
}

type AuthService struct {
	users  UserRepository
	tokens TokenGenerator
}

func NewAuthService(
	users UserRepository,
	tokens TokenGenerator,
) *AuthService {
	return &AuthService{
		users:  users,
		tokens: tokens,
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

func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (dto.LoginData, error) {
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.LoginData{},
				apperror.Unauthorized("用户名或密码错误")
		}

		return dto.LoginData{}, apperror.Internal(err)
	}

	if err := password.Verify(user.PasswordHash, req.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return dto.LoginData{},
				apperror.Unauthorized("用户名或密码错误")
		}

		return dto.LoginData{}, apperror.Internal(err)
	}

	accessToken, expiresIn, err := s.tokens.Generate(user)
	if err != nil {
		return dto.LoginData{}, apperror.Internal(err)
	}

	return dto.LoginData{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
		},
	}, nil
}
