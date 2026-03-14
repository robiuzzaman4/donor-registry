package response

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry/internal/domain"
)

// Pagination represents pagination metadata
type Pagination struct {
	Total           int64 `json:"total" db:"total"`
	Page            int   `json:"page" db:"page"`
	Limit           int   `json:"limit" db:"limit"`
	TotalPages      int   `json:"totalPages" db:"total_pages"`
	HasNextPage     bool  `json:"hasNextPage" db:"has_next_page"`
	HasPreviousPage bool  `json:"hasPreviousPage" db:"has_previous_page"`
	NextPage        *int  `json:"nextPage" db:"next_page"`
	PrevPage        *int  `json:"prevPage" db:"prev_page"`
}

// APIResponse represents standard API response structure
type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// BuildPagination creates pagination metadata
func BuildPagination(total int64, page, limit int) *Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNextPage := page < totalPages
	hasPreviousPage := page > 1

	var nextPage, prevPage *int
	if hasNextPage {
		next := page + 1
		nextPage = &next
	}
	if hasPreviousPage {
		prev := page - 1
		prevPage = &prev
	}

	return &Pagination{
		Total:           total,
		Page:            page,
		Limit:           limit,
		TotalPages:      totalPages,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		NextPage:        nextPage,
		PrevPage:        prevPage,
	}
}

// ParsePaginationParams extracts and validates pagination parameters from request
// Returns page, limit, and offset values with sensible defaults
func ParsePaginationParams(c *gin.Context) (page, limit, offset int) {
	// Default values
	page = 1
	limit = 10

	// Parse page parameter
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	// Parse limit parameter (max 100)
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Calculate offset
	offset = (page - 1) * limit

	return page, limit, offset
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Retrieved successfully",
		Data:    data,
	})
}

// SuccessWithMessage sends a successful response with custom message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination sends a successful response with pagination
func SuccessWithPagination(c *gin.Context, message string, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Created successfully",
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "Internal server error"

	// Map domain errors to HTTP status codes
	switch err {
	case domain.ErrNotFound, domain.ErrUserNotFound:
		statusCode = http.StatusNotFound
		message = err.Error()
	case domain.ErrInvalidCredentials, domain.ErrInvalidToken, domain.ErrTokenExpired:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case domain.ErrEmailExists, domain.ErrPhoneExists, domain.ErrAlreadyExists:
		statusCode = http.StatusConflict
		message = err.Error()
	case domain.ErrInvalidInput:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case domain.ErrForbidden:
		statusCode = http.StatusForbidden
		message = err.Error()
	default:
		if err != nil {
			message = err.Error()
		}
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

// NotFound sends a not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
