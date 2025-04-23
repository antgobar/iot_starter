package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateApiKey() string {
	return uuid.New().String()
}

func Encrypt(password string) (string, error) {
	p := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
