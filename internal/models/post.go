package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Title            string         `json:"title" gorm:"not null"`
	Content          string         `json:"content" gorm:"not null"`
	UserID           uint           `json:"user_id" gorm:"not null"`
	Category         string         `json:"category" gorm:"not null"`
	Tags             pq.StringArray `json:"tags" gorm:"type:text[]"`
	FeaturedImage    string         `json:"featured_image"`
	Status           string         `json:"status" gorm:"not null;default:draft"`
	ViewCount        uint           `json:"view_count" gorm:"not null;default:0"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Comments         []Comment
	LikesandDislikes []LikesandDislikes
}

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Comment   string         `json:"comment" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type LikesandDislikes struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	UserID       uint   `json:"user_id" gorm:"not null"`
	PostID       uint   `json:"post_id" gorm:"not null"`
	ReactionType string `json:"reaction_type" gorm:"not null"`
}
