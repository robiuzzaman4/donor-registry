package util

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/robiuzzaman4/donor-registry/internal/config"
	"github.com/robiuzzaman4/donor-registry/internal/domain"
)

// TokenClaims defines JWT claims used in this service
type TokenClaims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id,omitempty"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken func generate a new token
func GenerateToken(userID string, role string, duration time.Duration) (string, error) {
	claims := TokenClaims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	secret := config.GetConfig().JwtSecret
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ValidateToken validates a JWT and returns its claims
func ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	_ = ctx
	tokenString = strings.TrimSpace(tokenString)
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	}
	if tokenString == "" {
		return nil, domain.ErrInvalidToken
	}

	secret := config.GetConfig().JwtSecret
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrInvalidToken
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	if claims.ID == "" && claims.UserID != "" {
		claims.ID = claims.UserID
	}
	if claims.ID == "" {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
