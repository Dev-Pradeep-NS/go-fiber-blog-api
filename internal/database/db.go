package database

import (
	"github.com-Personal/go-fiber/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.LikesandDislikes{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
