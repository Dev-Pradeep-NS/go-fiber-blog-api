package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostHandler handles HTTP requests related to posts
type PostHandler struct {
	DB *gorm.DB
}

// NewPostHandler creates a new PostHandler instance
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

// GetPosts retrieves all posts from the database
func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	result := h.DB.Find(&posts)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch posts",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

// NewPost creates a new post in the database
func (h *PostHandler) NewPost(c *fiber.Ctx) error {
	newPost := new(models.Post)

	// Parse the request body
	if err := c.BodyParser(newPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}

	// Get the user ID from the context
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	newPost.UserID = userID

	// Validate post data
	if newPost.Title == "" || newPost.Content == "" || newPost.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}

	// Parse and validate tags
	tags := c.FormValue("tags")
	if tags == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one tag is required",
		})
	}
	newPost.Tags = strings.Split(tags, ",")

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	// Generate a unique filename for the image
	uniqueID := uuid.New()
	fileName := strings.Replace(uniqueID.String(), "-", "", -1)
	fileExt := strings.ToLower(filepath.Ext(file.Filename))

	// Validate file extension
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file type. Only JPG, JPEG, and PNG are allowed",
		})
	}

	image := fmt.Sprintf("%s%s", fileName, fileExt)
	uploadDir := "./uploads"
	uploadPath := fmt.Sprintf("%s/%s", uploadDir, image)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	newPost.FeaturedImage = image
	newPost.ViewCount = 0
	newPost.Status = "draft"

	// Use a transaction to ensure atomicity of database operations
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newPost).Error; err != nil {
			return err
		}

		if err := c.SaveFile(file, uploadPath); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Post",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newPost)
}

// GetPostById retrieves a single post by its ID
func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	// Fetch the post from the database
	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	// Increment the view count
	if err := h.DB.Model(&post).Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update view count",
			"error":   err.Error(),
		})
	}

	// Fetch the updated post with comments and likes/dislikes
	if err := h.DB.Preload("Comments").Preload("LikesandDislikes").First(&post, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch updated Post",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(post)
}

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	// Fetch the post from the database
	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	// Check if the user is authorized to update the post
	userID := c.Locals("user_id").(uint)
	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this post",
		})
	}

	// Parse the updated post data
	var updatedPost models.Post
	if err := c.BodyParser(&updatedPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}

	// Validate the updated post data
	if updatedPost.Title == "" || updatedPost.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}

	// Update the post in the database
	updateResult := h.DB.Model(&post).Omit("UserID", "ViewCount").Updates(updatedPost)
	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update Post",
			"error":   updateResult.Error.Error(),
		})
	}
	return c.JSON(post)
}

// DeletePost deletes a post from the database
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	// Fetch the post from the database
	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	// Check if the user is authorized to delete the post
	userID := c.Locals("user_id").(uint)
	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to delete this post",
		})
	}

	// Delete the post from the database
	deleteResult := h.DB.Delete(&post)
	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to delete Post",
			"error":   deleteResult.Error.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

// GetPostsByUser retrieves all posts for a specific user
func (h *PostHandler) GetPostsByUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var posts []models.Post
	result := h.DB.Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch posts",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}
