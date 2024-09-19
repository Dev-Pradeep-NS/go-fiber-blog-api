package middleware

import (
	"strings"

	"github.com-Personal/go-fiber/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		if accessToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing access token",
			})
		}

		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired access token",
			})
		}

		userID := uint(claims["user_id"].(float64))
		c.Locals("user_id", userID)
		c.Locals("username", claims["username"])
		return c.Next()
	}
}
