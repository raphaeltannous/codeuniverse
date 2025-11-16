package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateToken(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
