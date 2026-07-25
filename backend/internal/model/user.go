package model

import "time"

type UserRole string

const (
	UserRoleStudent UserRole = "student"
	UserRoleAdmin   UserRole = "admin"
)

type User struct {
	ID uint64 `gorm:"primaryKey" json:"id"`

	Username string `gorm:"type:varchar(32);not null;uniqueIndex" json:"username"`

	Name string `gorm:"type:varchar(32);not null" json:"name"`

	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`

	Role UserRole `gorm:"type:enum('student','admin');not null;default:'student';index" json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
