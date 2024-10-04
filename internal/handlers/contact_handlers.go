package handlers

import (
	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ContactHandlers struct {
	DB *gorm.DB
}

func NewContactHandlers(db *gorm.DB) *ContactHandlers {
	return &ContactHandlers{DB: db}
}

func (h *ContactHandlers) PostContact(c *fiber.Ctx) error {
	var newcontact models.Contact

	if err := c.BodyParser(&newcontact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "Unable to parse",
		})
	}

	if newcontact.Email == "" || newcontact.Message == "" || newcontact.Subject == "" || newcontact.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide required fields",
		})
	}

	result := h.DB.Create(&newcontact)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add contact details",
		})
	}
	return c.Status(fiber.StatusOK).JSON(newcontact)
}
