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

type ListPostsQuery struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Sort     string `form:"sort" binding:"oneof=latest"`
}

type PostListItemResponse struct {
	ID           uint64       `json:"id"`
	Content      string       `json:"content"`
	Author       UserResponse `json:"author"`
	LikeCount    int64        `json:"like_count"`
	CommentCount int64        `json:"comment_count"`
	CreatedAt    time.Time    `json:"created_at"`
}

type PaginationMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PostListResponse struct {
	Items []PostListItemResponse `json:"items"`
	Meta  PaginationMeta         `json:"meta"`
}

type CommentResponse struct {
	ID        uint64       `json:"id"`
	PostID    uint64       `json:"post_id"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

type PostDetailResponse struct {
	ID           uint64            `json:"id"`
	Content      string            `json:"content"`
	Author       UserResponse      `json:"author"`
	LikeCount    int64             `json:"like_count"`
	CommentCount int64             `json:"comment_count"`
	CreatedAt    time.Time         `json:"created_at"`
	Comments     []CommentResponse `json:"comments"`
}

type LikePostResponse struct {
	PostID  uint64 `json:"post_id"`
	IsLiked bool   `json:"is_liked"`
}

type GetLikeStatusesRequest struct {
	PostIDs []uint64 `json:"post_ids" binding:"required,min=1,max=100,dive,gt=0"`
}

type PostLikeStatus struct {
	PostID uint64 `json:"post_id"`
	Liked  bool   `json:"liked"`
}

type GetLikeStatusesResponse struct {
	Status []PostLikeStatus `json:"status"`
}
