package utils

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"git.riyt.dev/codeuniverse/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

func init() {
	devToken := "c3d64819480a"

	jwtSecret = os.Getenv("CODEUNIVERSE_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = devToken
		slog.Warn("jwtSecret is empty. Going to use devToken.", "devToken", devToken)
	} else {
		slog.Info("jwtSecret is updated.")
	}
}

func CreateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		slog.Error("failed to sign jwt", "err", err)
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("unexpected signing method")
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		slog.Warn("failed to validate jwt", "err", err)
	}

	if !token.Valid {
		slog.Warn("invalid JWT token")
		return nil, jwt.ErrTokenInvalidClaims
	}

	return token, nil
}
