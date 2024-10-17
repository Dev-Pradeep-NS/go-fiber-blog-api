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

func InitializeFirebaseApp() (*auth.Client, *storage.Client, error) {
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")
	config := &firebase.Config{
		StorageBucket: bucketName,
	}

	opt := option.WithCredentialsFile("serviceAccountKey.json")
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
