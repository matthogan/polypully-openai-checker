package secret

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

func GetValue(key string, ctx context.Context) (string, error) {

	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	credsBase64 := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	if credsBase64 == "" {
		return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS_JSON environment variable is not set")
	}

	creds, err := base64.StdEncoding.DecodeString(credsBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode credentials: %v", err)
	}

	client, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return "", fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
