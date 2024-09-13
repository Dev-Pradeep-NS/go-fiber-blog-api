package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Bio       string         `json:"bio"`
	AvatarURL string         `json:"avatar_url"`
	Followers []*User        `json:"followers" gorm:"many2many:user_followers;joinForeignKey:following_id;joinReferences:follower_id"`
	Following []*User        `json:"following" gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:following_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// BookmarkedPosts []*Post        `json:"bookmarks" gorm:"many2many:user_bookmarks;"`

	Posts            []Post             `json:"posts" gorm:"foreignKey:UserID"`
	Comments         []Comment          `json:"comments" gorm:"foreignKey:UserID"`
	LikesandDislikes []LikesandDislikes `json:"likesanddislikes" gorm:"foreignKey:UserID"`
	Bookmarks        []Bookmark         `json:"bookmarks" gorm:"foreignKey:UserID"`
}
