package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/robiuzzaman4/donor-registry/internal/config"
)

// GenerateToken func generate a new token
func GenerateToken(userID string, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.GetConfig().JwtSecret
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
