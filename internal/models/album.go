package models

type Album struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Title  string  `json:"title" gorm:"not null"`
	Artist string  `json:"artist" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
}
