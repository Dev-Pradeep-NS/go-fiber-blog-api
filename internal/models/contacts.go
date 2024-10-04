package models

type Contact struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email" gorm:"not null"`
	Subject string `json:"subject" gorm:"not null"`
	Message string `json:"message" gorm:"not null"`
}
