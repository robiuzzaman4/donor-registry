package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry/internal/domain"
	"github.com/robiuzzaman4/donor-registry/internal/rest/response"
	"github.com/robiuzzaman4/donor-registry/internal/util"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

// AuthGuard validates the access token and sets user data in context
func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractAccessToken(c)
		if tokenString == "" {
			response.Error(c, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := util.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		userID := claims.ID
		if userID == "" {
			response.Error(c, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		role := claims.Role
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, role)
		c.Next()
	}
}

func extractAccessToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
		return cookieToken
	}
	return ""
}
