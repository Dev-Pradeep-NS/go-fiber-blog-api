package models

type User struct {
	ID               uint               `json:"id" gorm:"primaryKey"`
	Username         string             `json:"username" gorm:"not null;unique"`
	Password         string             `json:"password" gorm:"not null"`
	Email            string             `json:"email" gorm:"not null;unique"`
	Posts            []Post             `json:"posts" gorm:"foreignKey:UserID"`
	Comments         []Comment          `json:"comments" gorm:"foreignKey:UserID"`
	LikesandDislikes []LikesandDislikes `json:"likesanddislikes" gorm:"foreignKey:UserID"`
}
