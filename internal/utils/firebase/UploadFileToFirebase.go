package firebase_utils

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com-Personal/go-fiber/config/firebase_config"
	"github.com-Personal/go-fiber/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadFileToFirebase(bucket *storage.BucketHandle, file io.Reader, uploadPath string) (string, error) {
	ctx := context.Background()
	obj := bucket.Object(uploadPath)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return attrs.MediaLink, nil
}

func UploadFileToFirebaseAndGetURL(c *fiber.Ctx, formFieldName, uploadDir string) (string, string, error) {
	fileHeader, err := c.FormFile(formFieldName)
	if err != nil {
		return "", "", err
	}

	uniqueID := uuid.New()
	fileName := strings.Replace(uniqueID.String(), "-", "", -1)
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))

	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		return "", "", fmt.Errorf("invalid file type")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	_, storageClient, err := firebase_config.InitializeFirebaseApp()
	if err != nil {
		return "", "", err
	}

	bucket, err := storageClient.Bucket(utils.GetSecretOrEnv("BUCKET_NAME"))
	if err != nil {
		return "", "", err
	}

	uploadPath := fmt.Sprintf("%s/%s%s", uploadDir, fileName, fileExt)

	imageURL, err := UploadFileToFirebase(bucket, file, uploadPath)
	if err != nil {
		return "", "", err
	}

	return imageURL, fileName + fileExt, nil
}
