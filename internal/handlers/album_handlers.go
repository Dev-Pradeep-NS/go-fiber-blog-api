package handlers

import (
	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AlbumHandler struct {
	DB *gorm.DB
}

func NewAlbumHandler(db *gorm.DB) *AlbumHandler {
	return &AlbumHandler{DB: db}
}

func (h *AlbumHandler) GetAlbums(c *fiber.Ctx) error {
	var albums []models.Album
	result := h.DB.Find(&albums)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch albums",
			"error":   result.Error.Error(),
		})
	}
	return c.JSON(albums)
}

func (h *AlbumHandler) CreateAlbum(c *fiber.Ctx) error {
	var newAlbum models.Album
	if err := c.BodyParser(&newAlbum); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse album data",
			"error":   err.Error(),
		})
	}

	if newAlbum.Title == "" || newAlbum.Artist == "" || newAlbum.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid album data",
		})
	}

	result := h.DB.Create(&newAlbum)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create album",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(newAlbum)
}

func (h *AlbumHandler) GetAlbumByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var album models.Album
	result := h.DB.First(&album, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Album not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch album",
			"error":   result.Error.Error(),
		})
	}
	return c.JSON(album)
}

func (h *AlbumHandler) UpdateAlbum(c *fiber.Ctx) error {
	id := c.Params("id")
	var album models.Album
	result := h.DB.First(&album, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Album not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to find album",
			"error":   result.Error.Error(),
		})
	}

	var updatedAlbum models.Album
	if err := c.BodyParser(&updatedAlbum); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse album data",
			"error":   err.Error(),
		})
	}

	if updatedAlbum.Title == "" || updatedAlbum.Artist == "" || updatedAlbum.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid album data",
		})
	}

	updateResult := h.DB.Model(&album).Updates(updatedAlbum)
	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update album",
			"error":   updateResult.Error.Error(),
		})
	}
	return c.JSON(album)
}

func (h *AlbumHandler) DeleteAlbum(c *fiber.Ctx) error {
	id := c.Params("id")
	var album models.Album
	result := h.DB.First(&album, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Album not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to find album",
			"error":   result.Error.Error(),
		})
	}

	deleteResult := h.DB.Delete(&album)
	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to delete album",
			"error":   deleteResult.Error.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Album deleted successfully",
	})
}
