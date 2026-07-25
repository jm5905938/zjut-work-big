package model

import "time"

type Comment struct {
	ID uint64 `gorm:"primaryKey" json:"id"`

	PostID uint64 `gorm:"not null;index" json:"post_id"`
	Post   Post   `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	AuthorID uint64 `gorm:"not null;index" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"author"`

	Content string `gorm:"type:text;not null" json:"content"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
