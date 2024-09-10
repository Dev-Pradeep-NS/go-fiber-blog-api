package handlers

import (
	"os"
	"time"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler struct holds the database connection
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// SafeUser struct represents user data safe for client-side exposure
type SafeUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var newUser models.User

	// Parse request body
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse user data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if newUser.Username == "" || newUser.Password == "" || newUser.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username and password are required",
		})
	}

	// Check for existing username
	var existingUser models.User
	if err := h.DB.Where("username = ?", newUser.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Check for existing email
	if err := h.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}
	newUser.Password = string(hashedPassword)

	// Create new user in the database
	result := h.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create new user",
			"error":   result.Error.Error(),
		})
	}

	// Prepare safe user data for response
	safeUser := SafeUser{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    safeUser,
	})
}

// Login handles user authentication
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse login data
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse login data",
			"error":   err.Error(),
		})
	}

	// Find user by email
	var user models.User
	result := h.DB.Where("email = ?", loginData.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// Generate JWT claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	t, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	// Prepare safe user data for response
	safeUser := SafeUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Return success response with token and user data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
		"user":    safeUser,
	})
}
