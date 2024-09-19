package handlers

import (
	"strconv"
	"time"

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

	parentID := c.Query("parent_id")
	if parentID != "" {
		parentIDUint, err := strconv.ParseUint(parentID, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid parent comment ID",
				"error":   err.Error(),
			})
		}
		parentIDUint32 := uint(parentIDUint)
		comment.ParentID = &parentIDUint32
	}

	userID := c.Locals("user_id").(uint)
	userName := c.Locals("username").(string)
	comment.UserID = userID
	comment.PostID = uint(num)
	comment.Username = userName
	comment.CreatedAt = time.Now()

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
func (h *CommentHandler) GetCommentsandCount(c *fiber.Ctx) error {
	postID := c.Params("id")

	var comments []models.Comment
	var count int64

	if err := h.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Replies.Replies.Replies.Replies").
		Find(&comments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch Comments",
		})
	}

	if err := h.DB.Model(&models.Comment{}).Where("post_id = ? AND parent_id is NULL", postID).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count Comments",
		})
	}

	return c.JSON(fiber.Map{
		"comments": comments,
		"count":    count,
	})
}
