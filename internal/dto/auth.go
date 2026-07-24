package dto

import (
	"github.com/jm5905938/zjut-work-big/internal/model"
)

type RegisterRequest struct {
	Username string         `json:"username" binding:"required,numeric,min=1,max=32"`
	Name     string         `json:"name" binding:"required,min=1,max=32"`
	Password string         `json:"password" binding:"required,min=8,max=16"`
	Role     model.UserRole `json:"role" binding:"required,oneof=student admin"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint64         `json:"id"`
	Username string         `json:"username"`
	Name     string         `json:"name"`
	Role     model.UserRole `json:"role"`
}

type LoginData struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}
