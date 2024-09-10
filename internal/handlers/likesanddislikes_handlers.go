package handlers

import (
	"strconv"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LikesandDislikes struct holds the database connection
type LikesandDislikes struct {
	DB *gorm.DB
}

// NewLikesandDislikes creates a new instance of LikesandDislikes
func NewLikesandDislikes(db *gorm.DB) *LikesandDislikes {
	return &LikesandDislikes{DB: db}
}

// AddReaction handles the addition or update of a reaction (like or dislike) to a post
func (h *LikesandDislikes) AddReaction(c *fiber.Ctx) error {
	postID := c.Params("id")
	var reaction models.LikesandDislikes

	// Parse the request body into the reaction struct
	if err := c.BodyParser(&reaction); err != nil {
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

	// Set the user ID and post ID for the reaction
	userID := c.Locals("user_id").(uint)
	reaction.UserID = userID
	reaction.PostID = uint(num)

	// Check if a reaction already exists for this user and post
	var existingReaction models.LikesandDislikes
	result := h.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingReaction)
	if result.Error == nil {
		// Update existing reaction
		existingReaction.ReactionType = reaction.ReactionType
		updateResult := h.DB.Save(&existingReaction)
		if updateResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update reaction",
				"error":   updateResult.Error.Error(),
			})
		}
		return c.JSON(existingReaction)
	}

	// Create new reaction if it doesn't exist
	createResult := h.DB.Create(&reaction)
	if createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add reaction",
			"error":   createResult.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(reaction)
}

// RemoveReaction handles the removal of a reaction (like or dislike) from a post
func (h *LikesandDislikes) RemoveReaction(c *fiber.Ctx) error {
	postID := c.Params("id")
	userID := c.Locals("user_id").(uint)

	// Find the existing reaction
	var reaction models.LikesandDislikes
	result := h.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&reaction)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reaction not found",
		})
	}

	// Delete the reaction
	deleteResult := h.DB.Delete(&reaction)
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
