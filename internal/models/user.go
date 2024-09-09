package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Posts    []Post `json:"posts" gorm:"foreignKey:UserID"`
}
