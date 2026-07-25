package dto

import "time"

type CreatePostRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}

type PostResponse struct {
	ID        uint64       `json:"id"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}
