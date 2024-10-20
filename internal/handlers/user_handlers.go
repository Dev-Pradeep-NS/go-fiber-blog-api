package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com-Personal/go-fiber/internal/models"
	"github.com-Personal/go-fiber/internal/utils"
	firebase_utils "github.com-Personal/go-fiber/internal/utils/firebase"
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

	accessToken, err := generateToken(newUser.ID, newUser.Username, time.Hour)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
	}

	refreshToken, err := generateToken(newUser.ID, newUser.Username, time.Hour*24*7)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate refresh token",
			"error":   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   false,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Strict",
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "User registered successfully",
		"user":          safeUser,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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

	accessToken, err := generateToken(user.ID, user.Username, time.Hour)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
	}

	refreshToken, err := generateToken(user.ID, user.Username, time.Hour*24*7)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate refresh token",
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

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   false,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Strict",
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          safeUser,
	})
}

func (h *UserHandler) CheckEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	fmt.Println(email)

	var user models.User
	result := h.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Email not found",
			})
		}
		// Handle any other database errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email found",
	})
}

func (h *UserHandler) GetAllUsernameAndEmails(c *fiber.Ctx) error {

	var users []struct {
		Email    string
		Username string
	}
	if result := h.DB.Model(models.User{}).Select("email, username").Find(&users); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to fetch emails and usernames",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Users": users,
	})
}

func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	var data struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse request body",
			"error":   err.Error(),
		})
	}

	// Validate input
	if data.Email == "" || data.OldPassword == "" || data.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email, old password, and new password are required",
		})
	}

	var user models.User
	result := h.DB.Where("email = ?", data.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   result.Error.Error(),
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid old password",
		})
	}

	// Check if new password is different from old password
	if data.OldPassword == data.NewPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "New password must be different from old password",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	user.Password = string(hashedPassword)
	if err := h.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successful",
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   false,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return c.SendStatus(fiber.StatusOK)
}

func generateToken(userID uint, username string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(utils.GetSecretOrEnv("JWT_SECRET_KEY"))
	return token.SignedString(secretKey)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	log.Printf("Received refresh token: %s", refreshToken)
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing refresh token",
			"error":   "No refresh token",
		})
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired refresh token",
		})
	}

	userID := uint(claims["user_id"].(float64))
	username := claims["username"].(string)

	newAccessToken, err := generateToken(userID, username, time.Hour)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate new access token",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

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

func (h *UserHandler) GetUserDetail(c *fiber.Ctx) error {
	username := c.Params("username")
	var user models.User
	if err := h.DB.Where("username = ?", username).Find(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
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

func (h *UserHandler) GetAvatarImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filepath := filepath.Join("./uploads/avatars", filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	return c.SendFile(filepath)
}

func (h *UserHandler) UploadAvatar(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	imageURL, _, err := firebase_utils.UploadFileToFirebaseAndGetURL(c, "avatar", "avatars")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload avatar",
			"error":   err.Error(),
		})
	}

	user.AvatarURL = imageURL

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
	followerID := c.Locals("user_id").(uint)
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
	followerID := c.Locals("user_id").(uint)
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
