package main

import (
	"fmt"
	"log"

	"github.com-Personal/go-fiber/config"
	"github.com-Personal/go-fiber/internal/database"
	"github.com-Personal/go-fiber/internal/handlers"
	"github.com-Personal/go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to the database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize Firebase
	_, _, err = config.InitializeFirebaseApp()
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err) // Ensure the app terminates on error
	}
	fmt.Println("Firebase Auth client initialized successfully.")

	// Initialize Fiber router
	router := fiber.New()
	router.Use(middleware.CorsMiddleware())
	router.Use(logger.New())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	likes_and_dislikes := handlers.NewLikesandDislikes(db)
	bookmarkHandler := handlers.NewBookmarkHandler(db)
	contactHandler := handlers.NewContactHandlers(db)

	// Public routes
	router.Post("/login", userHandler.Login)
	router.Post("/register", userHandler.Register)
	router.Post("/refresh", userHandler.RefreshToken)
	router.Post("/logout", userHandler.Logout)
	router.Get("/verifyemail/:email", userHandler.CheckEmail)
	router.Put("/reset-password", userHandler.ForgotPassword)

	// Protected routes group
	api := router.Group("/", middleware.AuthMiddleware())

	// User routes
	users := api.Group("/users")
	users.Get("", userHandler.GetProfile)
	users.Get("/:username", userHandler.GetUserDetail)
	users.Get("/uploads/avatars/:filename", userHandler.GetAvatarImage)
	users.Put("/:id", userHandler.UpdateProfile)
	users.Post("/:id/avatar", userHandler.UploadAvatar)
	users.Post("/follow/:followingID", userHandler.FollowUser)
	users.Delete("/unfollow/:followingID", userHandler.UnfollowUser)
	users.Get("/:id/followers", userHandler.GetFollowers)
	users.Get("/:id/following", userHandler.GetFollowing)
	api.Get("/users-emails", userHandler.GetAllUsernameAndEmails)

	// Reaction routes
	api.Post("/posts/:id/like", likes_and_dislikes.LikePost)
	api.Post("/posts/:id/dislike", likes_and_dislikes.DisLikePost)
	api.Get("/posts/:post_id/reactions", likes_and_dislikes.GetReaction)

	// Comment routes
	api.Post("/posts/:id/comments", commentHandler.AddComment)
	api.Get("/posts/:id/comments", commentHandler.GetCommentsandCount)
	api.Put("/comments/:id", commentHandler.UpdateComment)
	api.Delete("/comments/:id", commentHandler.DeleteComment)

	// Post routes
	api.Get("/posts", postHandler.GetPosts)
	api.Get("/uploads/:filename", postHandler.GetImage)
	api.Post("/posts", postHandler.NewPost)
	api.Get("posts/:username/:slug", postHandler.GetPostBySlug)
	api.Put("/posts/:id", postHandler.UpdatePost)
	api.Delete("/posts/:id", postHandler.DeletePost)
	api.Get("/users/:id/posts", postHandler.GetPostsByUser)

	api.Post("/users/:post_id/bookmark", bookmarkHandler.BookmarkPost)
	api.Get("/users/post/bookmarks", bookmarkHandler.GetBookmarks)
	api.Get("/:post_id/bookmarkscount", bookmarkHandler.GetBookmarkCount)

	// Contact routes
	api.Post("/contact-us", contactHandler.PostContact)

	// Start the server
	log.Fatal(router.Listen(cfg.HOST + ":" + cfg.PORT))
}
