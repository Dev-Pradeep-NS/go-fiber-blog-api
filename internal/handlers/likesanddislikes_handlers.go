package handlers

import (
	"strconv"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LikesandDislikes struct {
	DB *gorm.DB
}

func NewLikesandDislikes(db *gorm.DB) *LikesandDislikes {
	return &LikesandDislikes{DB: db}
}

func (h *LikesandDislikes) GetReaction(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	var reactions []models.LikesandDislikes

	result := h.DB.Where("post_id = ?", postID).Find(&reactions)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch reactions",
		})
	}

	if len(reactions) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "No reactions found for this post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(reactions)
}

func (h *LikesandDislikes) LikePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	var reaction models.LikesandDislikes

	if err := c.BodyParser(&reaction); err != nil {
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
	reaction.UserID = userID
	reaction.PostID = uint(num)
	reaction.ReactionType = "like"

	var existingReaction models.LikesandDislikes

	result := h.DB.Where("user_id = ? AND post_id = ? AND reaction_type = ?", userID, num, "like").First(&existingReaction)
	if result.RowsAffected > 0 {
		deleteResult := h.DB.Delete(&existingReaction)
		if deleteResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unable to remove reaction",
				"error":   deleteResult.Error.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Reaction removed successfully",
		})
	}

	dislikeResult := h.DB.Where("user_id = ? AND post_id = ? AND reaction_type = ?", userID, num, "dislike").First(&existingReaction)
	if dislikeResult.RowsAffected > 0 {
		deleteResult := h.DB.Delete(&existingReaction)
		if deleteResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unable to remove reaction",
				"error":   deleteResult.Error.Error(),
			})
		}
	}

	createResult := h.DB.Create(&reaction)
	if createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add reaction",
			"error":   createResult.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(reaction)
}

func (h *LikesandDislikes) DisLikePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	var reaction models.LikesandDislikes

	if err := c.BodyParser(&reaction); err != nil {
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
	reaction.UserID = userID
	reaction.PostID = uint(num)
	reaction.ReactionType = "dislike"

	var existingReaction models.LikesandDislikes

	result := h.DB.Where("user_id = ? AND post_id = ? AND reaction_type = ?", userID, num, "dislike").First(&existingReaction)
	if result.RowsAffected > 0 {
		deleteResult := h.DB.Delete(&existingReaction)
		if deleteResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unable to remove reaction",
				"error":   deleteResult.Error.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Reaction removed successfully",
		})
	}

	likeResult := h.DB.Where("user_id = ? AND post_id = ? AND reaction_type = ?", userID, num, "like").First(&existingReaction)
	if likeResult.RowsAffected > 0 {
		deleteResult := h.DB.Delete(&existingReaction)
		if deleteResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unable to remove reaction",
				"error":   deleteResult.Error.Error(),
			})
		}
	}

	createResult := h.DB.Create(&reaction)
	if createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add reaction",
			"error":   createResult.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(reaction)
}
