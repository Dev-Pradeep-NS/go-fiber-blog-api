package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a Fiber middleware for JWT authentication
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		// Remove the "Bearer " prefix from the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// Ensure the token uses the correct signing method (HS256)
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token signing method",
			})
		}

		// Check if the token is valid and not expired
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired token",
			})
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token claims",
			})
		}

		// Extract and validate the user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid user ID in token",
			})
		}
		userID := uint(userIDFloat)

		// Store user information in the context for later use
		c.Locals("user_id", userID)
		c.Locals("username", claims["username"])

		// Proceed to the next middleware or route handler
		return c.Next()
	}
}
