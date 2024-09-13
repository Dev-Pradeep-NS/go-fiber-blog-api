package main

import (
	"fmt"
	"log"

	"github.com-Personal/go-fiber/config"
	"github.com-Personal/go-fiber/internal/database"
	"github.com-Personal/go-fiber/internal/handlers"
	"github.com-Personal/go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Println(cfg.SERVER_URL)

	// Connect to the database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize Fiber router
	router := fiber.New()
	router.Use(middleware.CorsMiddleware())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	likes_and_dislikes := handlers.NewLikesandDislikes(db)
	bookmarkHandler := handlers.NewBookmarkHandler(db)

	// Public routes
	router.Post("/login", userHandler.Login)
	router.Post("/register", userHandler.Register)

	// Protected routes group
	api := router.Group("/", middleware.AuthMiddleware())

	// User routes
	users := api.Group("/users")
	users.Get("/", userHandler.GetProfile)
	users.Get("/uploads/avatars/:filename", userHandler.GetAvatarImage)
	users.Put("/:id", userHandler.UpdateProfile)
	users.Post("/:id/avatar", userHandler.UploadAvatar)
	users.Post("/:followerID/follow/:followingID", userHandler.FollowUser)
	users.Delete("/:followerID/unfollow/:followingID", userHandler.UnfollowUser)
	users.Get("/:id/followers", userHandler.GetFollowers)
	users.Get("/:id/following", userHandler.GetFollowing)

	// Post routes
	api.Get("/posts", postHandler.GetPosts)
	api.Get("/uploads/:filename", postHandler.GetImage)
	api.Post("/posts", postHandler.NewPost)
	api.Get("posts/:username/:slug", postHandler.GetPostBySlug)
	api.Put("/posts/:id", postHandler.UpdatePost)
	api.Delete("/posts/:id", postHandler.DeletePost)
	api.Get("/users/:id/posts", postHandler.GetPostsByUser)

	api.Post("/users/:post_id/bookmark", bookmarkHandler.BookmarkPost)
	api.Get("/:post_id/bookmarkscount", bookmarkHandler.GetBookmarkCount)

	// Comment routes
	api.Post("/posts/:id/comments", commentHandler.AddComment)
	api.Put("/comments/:id", commentHandler.UpdateComment)
	api.Delete("/comments/:id", commentHandler.DeleteComment)

	// Reaction routes
	api.Post("/posts/:id/reactions", likes_and_dislikes.AddReaction)
	api.Delete("/posts/:id/reactions", likes_and_dislikes.RemoveReaction)

	// Start the server
	log.Fatal(router.Listen(cfg.SERVER_URL))
}
