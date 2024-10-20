package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func GetSecret(key, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening secrets file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split each line into key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip if the line is not a valid key-value pair
		}

		currentKey := strings.TrimSpace(parts[0])
		currentValue := strings.TrimSpace(parts[1])

		// Return the value for the specified key
		if currentKey == key {
			return currentValue, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading secrets file: %w", err)
	}

	return "", errors.New("key not found: " + key)
}

func GetSecretOrEnv(key string) string {
	secretValue, err := GetSecret(key, "/run/secrets/"+key)
	if err == nil {
		return secretValue
	}

	return os.Getenv(key)
}
