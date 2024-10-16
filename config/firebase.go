package config

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

// InitializeFirebaseApp initializes the Firebase app and returns the Auth client and Storage client
func InitializeFirebaseApp() (*auth.Client, *storage.Client, error) {
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME") // Ensure this env variable is set
	config := &firebase.Config{
		StorageBucket: bucketName,
	}

	// Load the service account key
	opt := option.WithCredentialsFile("serviceAccountKey.json") // Path to your Firebase service account file
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	// Initialize Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	// Initialize Firebase Storage client
	storageClient, err := app.Storage(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing storage client: %v", err)
	}

	return authClient, storageClient, nil
}
