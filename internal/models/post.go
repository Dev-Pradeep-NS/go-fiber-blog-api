package models

import "time"

type Post struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title" gorm:"not null" validate:"required,min=1,max=255"`
	Content   string     `json:"content" gorm:"not null" validate:"required"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
