package firebase_config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"github.com-Personal/go-fiber/internal/utils"
	"google.golang.org/api/option"
)

func InitializeFirebaseApp() (*auth.Client, *storage.Client, error) {
	ctx := context.Background()
	bucketName := utils.GetSecretOrEnv("BUCKET_NAME")
	config := &firebase.Config{
		StorageBucket: bucketName,
	}

	creds := map[string]interface{}{
		"type":                        utils.GetSecretOrEnv("FIREBASE_TYPE"),
		"project_id":                  utils.GetSecretOrEnv("FIREBASE_PROJECT_ID"),
		"private_key_id":              utils.GetSecretOrEnv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 strings.Replace(utils.GetSecretOrEnv("FIREBASE_PRIVATE_KEY"), "\\n", "\n", -1),
		"client_email":                utils.GetSecretOrEnv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   utils.GetSecretOrEnv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    utils.GetSecretOrEnv("FIREBASE_AUTH_URI"),
		"token_uri":                   utils.GetSecretOrEnv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": utils.GetSecretOrEnv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        utils.GetSecretOrEnv("FIREBASE_CLIENT_X509_CERT_URL"),
		"universe_domain":             utils.GetSecretOrEnv("FIREBASE_UNIVERSE_DOMAIN"),
	}

	// Convert the creds map to a JSON byte array
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		log.Printf("error marshalling Firebase credentials: %v", err)
		return nil, nil, fmt.Errorf("error marshalling Firebase credentials: %v", err)
	}

	// Use the JSON byte array with the Firebase option
	opt := option.WithCredentialsJSON(credsJSON)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	storageClient, err := app.Storage(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing storage client: %v", err)
	}

	return authClient, storageClient, nil
}
