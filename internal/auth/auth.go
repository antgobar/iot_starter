package auth

import "github.com/google/uuid"

func GenerateApiKey() string {
	return uuid.New().String()
}
