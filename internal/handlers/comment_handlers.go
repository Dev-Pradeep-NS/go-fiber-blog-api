package handlers

import (
	"strconv"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CommentHandler handles HTTP requests related to comments
type CommentHandler struct {
	DB *gorm.DB
}

// NewCommentHandler creates a new CommentHandler instance
func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

// AddComment handles the creation of a new comment
func (h *CommentHandler) AddComment(c *fiber.Ctx) error {
	postID := c.Params("id")
	var comment models.Comment

	// Parse the request body into the comment struct
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}

	// Convert the post ID from string to uint
	num, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post ID",
			"error":   err.Error(),
		})
	}

	// Set the user ID and post ID for the comment
	userID := c.Locals("user_id").(uint)
	comment.UserID = userID
	comment.PostID = uint(num)

	// Save the comment to the database
	result := h.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add comment",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

// UpdateComment handles the updating of an existing comment
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment

	// Find the comment by ID
	result := h.DB.First(&comment, commentID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}

	// Check if the user is authorized to update the comment
	userID := c.Locals("user_id").(uint)
	if comment.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this comment",
		})
	}

	// Parse the updated comment data from the request body
	var updatedComment models.Comment
	if err := c.BodyParser(&updatedComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}

	// Update the comment content
	comment.Comment = updatedComment.Comment
	updateResult := h.DB.Save(&comment)
	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update comment",
			"error":   updateResult.Error.Error(),
		})
	}

	return c.JSON(comment)
}

// DeleteComment handles the deletion of an existing comment
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment

	// Find the comment by ID
	result := h.DB.First(&comment, commentID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}

	// Check if the user is authorized to delete the comment
	userID := c.Locals("user_id").(uint)
	if comment.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to delete this comment",
		})
	}

	// Delete the comment from the database
	deleteResult := h.DB.Delete(&comment)
	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to delete comment",
			"error":   deleteResult.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}
