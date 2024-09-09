package handlers

import (
	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	result := h.DB.Find(&posts)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "failed to fetch posts",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(posts)
}

func (h *PostHandler) NewPost(c *fiber.Ctx) error {
	var newPost models.Post
	if err := c.BodyParser(&newPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
		})
	}

	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized or missing user ID",
		})
	}

	userID := uint(userIDFloat)
	newPost.UserID = userID

	if newPost.Title == "" || newPost.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}
	result := h.DB.Create(&newPost)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Post",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(newPost)
}

func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusFound).JSON(post)
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized or missing user ID",
		})
	}

	userID := uint(userIDFloat)
	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this post",
		})
	}

	var updatedPost models.Post
	if err := c.BodyParser(&updatedPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
		})
	}

	if updatedPost.Title == "" || updatedPost.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}

	updateResult := h.DB.Model(&post).Omit("UserID").Updates(updatedPost)
	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update Post",
			"error":   updateResult.Error.Error(),
		})
	}
	return c.JSON(post)
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized or missing user ID",
		})
	}

	userID := uint(userIDFloat)

	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this post",
		})
	}

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
