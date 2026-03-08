package domain

import "errors"

var (
	// Common errors
	ErrNotFound       = errors.New("Resource not found")
	ErrAlreadyExists  = errors.New("Resource already exists")
	ErrInvalidInput   = errors.New("Invalid input")
	ErrUnauthorized   = errors.New("Unauthorized")
	ErrForbidden      = errors.New("Forbidden")
	ErrInternalServer = errors.New("Internal server error")

	// Auth errors
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrInvalidToken       = errors.New("Invalid token")
	ErrTokenExpired       = errors.New("Token expired")
	ErrEmailExists        = errors.New("Email already exists")
	ErrPhoneExists        = errors.New("Phone number already exists")

	// User errors
	ErrUserNotFound = errors.New("User not found")
)
