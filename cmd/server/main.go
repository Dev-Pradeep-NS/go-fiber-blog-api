package main

import (
	"log"

	"github.com-Personal/go-fiber/config"
	"github.com-Personal/go-fiber/internal/database"
	"github.com-Personal/go-fiber/internal/handlers"
	"github.com-Personal/go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	router := fiber.New()
	router.Use(middleware.CorsMiddleware())

	albumHandler := handlers.NewAlbumHandler(db)
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)

	router.Post("/login", userHandler.Login)
	router.Post("/register", userHandler.Register)

	api := router.Group("/", middleware.AuthMiddleware())

	api.Get("albums", albumHandler.GetAlbums)
	api.Post("albums", albumHandler.CreateAlbum)
	api.Get("albums/:id", albumHandler.GetAlbumByID)
	api.Put("albums/:id", albumHandler.UpdateAlbum)
	api.Delete("albums/:id", albumHandler.DeleteAlbum)

	api.Get("posts", postHandler.GetPosts)
	api.Post("posts", postHandler.NewPost)
	api.Get("posts/:id", postHandler.GetPostById)
	api.Put("posts/:id", postHandler.UpdatePost)
	api.Delete("posts/:id", postHandler.DeletePost)

	log.Fatal(router.Listen(":" + cfg.PORT))
}
