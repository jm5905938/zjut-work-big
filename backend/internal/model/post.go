package model

import "time"

type Post struct {
	ID uint64 `gorm:"primaryKey" json:"id"`

	AuthorID uint64 `gorm:"not null;index" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"author"`

	Content string `gorm:"type:text;not null" json:"content"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LikeCount    int64     `gorm:"->;-:migration" json:"like_count"`
	CommentCount int64     `gorm:"->;-:migration" json:"comment_count"`
	Comments     []Comment `gorm:"foreignKey:PostID" json:"comments"`
}
