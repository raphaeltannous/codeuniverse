package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"
)

var signingSecret []byte

func init() {
	signingSecret = make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, signingSecret)
	if err != nil {
		log.Fatal("failed to read random string.")
	}

}

func GenerateSignedUrl(urlPath string, expiresInNMinutes int) string {
	expiresAt := time.Now().Add(time.Duration(expiresInNMinutes) * time.Minute).Unix()

	payload := fmt.Sprintf("%s:%d", urlPath, expiresAt)

	h := hmac.New(sha256.New, signingSecret)
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s?expires=%d&signature=%s", urlPath, expiresAt, signature)
}

func ValidateSignedUrl(urlPath, signature string, expiresAt int) bool {
	payload := fmt.Sprintf("%s:%d", urlPath, expiresAt)

	h := hmac.New(sha256.New, signingSecret)
	h.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
