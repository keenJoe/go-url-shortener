package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"` // json:"-" 表示不输出到JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}
