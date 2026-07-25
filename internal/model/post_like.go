package model

import "time"

type PostLike struct {
	PostID uint64 `gorm:"primaryKey" json:"post_id"`
	Post   Post   `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	UserID uint64 `gorm:"primaryKey;index" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"created_at"`
}
