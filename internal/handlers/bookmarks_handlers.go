package handlers

import (
	"strconv"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BookmarkHandler struct {
	DB *gorm.DB
}

func NewBookmarkHandler(db *gorm.DB) *BookmarkHandler {
	return &BookmarkHandler{DB: db}
}

func (h *BookmarkHandler) BookmarkPost(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized user",
		})
	}

	var bookmark models.Bookmark

	var count int64
	h.DB.Where("user_id = ?", userID).Where("post_id = ?", postID).Find(&bookmark).Count(&count)
	if count > 0 {
		result := h.DB.Delete(&bookmark)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unable to remove bookmark",
			})
		}
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Bookmark removed successfully",
		})
	}

	postIDUint, err := strconv.ParseUint(postID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}
	bookmark.PostID = uint(postIDUint)
	bookmark.UserID = uint(userID)

	result := h.DB.Create(&bookmark)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to bookmark the post",
		})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":         "Bookmarked successfully",
		"bookmarked post": bookmark,
	})
}

func (h *BookmarkHandler) GetBookmarkCount(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	var post models.Post
	findPost := h.DB.First(&post, postID)
	if findPost.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not available",
		})
	}

	var count int64

	result := h.DB.Model(&models.Bookmark{}).Where("post_id = ?", postID).Count(&count)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"BookmarkedPostCount": count,
	})
}
