package handlers

import (
	"strconv"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

func (h *CommentHandler) AddComment(c *fiber.Ctx) error {
	postID := c.Params("id")
	var comment models.Comment

	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}

	num, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post ID",
			"error":   err.Error(),
		})
	}

	userID := c.Locals("user_id").(uint)
	comment.UserID = userID
	comment.PostID = uint(num)

	result := h.DB.Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add comment",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment
	result := h.DB.First(&comment, commentID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}
	userID := c.Locals("user_id").(uint)
	if comment.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this comment",
		})
	}
	var updatedComment models.Comment
	if err := c.BodyParser(&updatedComment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}
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

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment
	result := h.DB.First(&comment, commentID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}
	userID := c.Locals("user_id").(uint)
	if comment.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to delete this comment",
		})
	}
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
