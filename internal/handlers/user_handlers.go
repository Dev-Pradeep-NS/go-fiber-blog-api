package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com-Personal/go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type SafeUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var userReg UserRegistration

	if err := c.BodyParser(&userReg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse user data",
			"error":   err.Error(),
		})
	}

	log.Printf("received : %+v", userReg)

	if userReg.Username == "" || userReg.Password == "" || userReg.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username, password, and email are required",
		})
	}

	var existingUser models.User
	if err := h.DB.Where("username = ?", userReg.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	if err := h.DB.Where("email = ?", userReg.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	newUser := models.User{
		Username: userReg.Username,
		Email:    userReg.Email,
		Password: string(hashedPassword),
	}

	result := h.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create new user",
			"error":   result.Error.Error(),
		})
	}

	safeUser := SafeUser{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Bio:       newUser.Bio,
		AvatarURL: newUser.AvatarURL,
		CreatedAt: newUser.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    safeUser,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse login data",
			"error":   err.Error(),
		})
	}

	var user models.User
	result := h.DB.Where("email = ?", loginData.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	t, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	safeUser := SafeUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
		"user":    safeUser,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	safeUser := SafeUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(safeUser)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse update data",
			"error":   err.Error(),
		})
	}

	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Bio != "" {
		user.Bio = updateData.Bio
	}

	if err := h.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
	}

	safeUser := SafeUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    safeUser,
	})
}

func (h *UserHandler) UploadAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Ensure the upload directory exists
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	// Save the file
	if err := c.SaveFile(file, fmt.Sprintf("%s/%s", uploadDir, filename)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}

	// Update the user's avatar URL in the database
	user.AvatarURL = fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := h.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update avatar URL",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Avatar uploaded successfully",
		"avatarURL": user.AvatarURL,
	})
}

func (h *UserHandler) FollowUser(c *fiber.Ctx) error {
	followerID := c.Params("followerID")
	followingID := c.Params("followingID")

	var follower, following models.User
	if err := h.DB.First(&follower, followerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Follower not found",
		})
	}
	if err := h.DB.First(&following, followingID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User to follow not found",
		})
	}

	if err := h.DB.Model(&follower).Association("Following").Append(&following); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to follow user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully followed user",
	})
}

func (h *UserHandler) UnfollowUser(c *fiber.Ctx) error {
	followerID := c.Params("followerID")
	followingID := c.Params("followingID")

	var follower, following models.User
	if err := h.DB.First(&follower, followerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Follower not found",
		})
	}
	if err := h.DB.First(&following, followingID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User to unfollow not found",
		})
	}

	if err := h.DB.Model(&follower).Association("Following").Delete(&following); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to unfollow user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully unfollowed user",
	})
}

func (h *UserHandler) GetFollowers(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.Preload("Followers").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	var safeFollowers []SafeUser
	for _, follower := range user.Followers {
		safeFollowers = append(safeFollowers, SafeUser{
			ID:        follower.ID,
			Username:  follower.Username,
			Email:     follower.Email,
			Bio:       follower.Bio,
			AvatarURL: follower.AvatarURL,
			CreatedAt: follower.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"followers": safeFollowers,
	})
}

func (h *UserHandler) GetFollowing(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.Preload("Following").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	var safeFollowing []SafeUser
	for _, following := range user.Following {
		safeFollowing = append(safeFollowing, SafeUser{
			ID:        following.ID,
			Username:  following.Username,
			Email:     following.Email,
			Bio:       following.Bio,
			AvatarURL: following.AvatarURL,
			CreatedAt: following.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"following": safeFollowing,
	})
}
