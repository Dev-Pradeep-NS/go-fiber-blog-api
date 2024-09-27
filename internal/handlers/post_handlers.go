package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com-Personal/go-fiber/internal/models"
	"github.com-Personal/go-fiber/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

func (h *PostHandler) GetImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filepath := filepath.Join("./uploads", filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	return c.SendFile(filepath)
}

func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	result := h.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "bio", "avatar_url", "created_at")
		}).
		Preload("Comments").
		Preload("LikesandDislikes").
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch posts",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

func (h *PostHandler) NewPost(c *fiber.Ctx) error {
	newPost := new(models.Post)

	if err := c.BodyParser(newPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}
	newPost.UserID = userID
	if newPost.Title == "" || newPost.Content == "" || newPost.Category == "" || newPost.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}

	tags := c.FormValue("tags")
	if tags == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one tag is required",
		})
	}
	newPost.Tags = strings.Split(tags, ",")

	newPost.Slug = utils.CreateSlug(newPost.Title)

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	uniqueID := uuid.New()
	fileName := strings.Replace(uniqueID.String(), "-", "", -1)
	fileExt := strings.ToLower(filepath.Ext(file.Filename))

	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file type. Only JPG, JPEG, and PNG are allowed",
		})
	}

	image := fmt.Sprintf("%s%s", fileName, fileExt)
	uploadDir := "./uploads"
	uploadPath := fmt.Sprintf("%s/%s", uploadDir, image)
	imageUrl := fmt.Sprintf("http://localhost:8000/uploads/%s", image)

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	newPost.FeaturedImage = image
	newPost.FeaturedImageUrl = imageUrl
	newPost.ViewCount = 0
	newPost.Status = "draft"

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newPost).Error; err != nil {
			return err
		}

		if err := c.SaveFile(file, uploadPath); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create Post",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newPost)
}

func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	userID := c.Locals("user_id").(uint)
	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to update this post",
		})
	}

	var updatedPost models.Post
	if err := c.BodyParser(&updatedPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse the data",
			"error":   err.Error(),
		})
	}
	if updatedPost.Title == "" || updatedPost.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Post data",
		})
	}

	updateResult := h.DB.Model(&post).Omit("UserID", "ViewCount").Updates(updatedPost)
	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update Post",
			"error":   updateResult.Error.Error(),
		})
	}
	return c.JSON(post)
}

func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post

	result := h.DB.First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch Post",
			"error":   result.Error.Error(),
		})
	}

	userID := c.Locals("user_id").(uint)
	if post.UserID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to delete this post",
		})
	}

	deleteResult := h.DB.Delete(&post)
	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to delete Post",
			"error":   deleteResult.Error.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}

func (h *PostHandler) GetPostsByUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var posts []models.Post
	result := h.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email")
		}).
		Where("user_id = ?", userID).Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch posts",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

func (h *PostHandler) GetPostBySlug(c *fiber.Ctx) error {
	username := c.Params("username")
	slug := c.Params("slug")

	var post models.Post
	postResult := h.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email")
		}).
		Preload("User").
		Joins("JOIN users ON posts.user_id = users.id").
		Where("users.username = ? AND posts.slug = ?", username, slug).
		First(&post)

	if postResult.Error != nil {
		if postResult.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Post not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve post",
		})
	}

	if err := h.DB.Model(&post).Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update view count",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(post)
}
