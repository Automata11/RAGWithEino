package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ragwitheitno/internal/auth"
)

// AuthRequired middleware validates JWT token
func AuthRequired(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminRequired middleware checks if user has admin role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserID gets the current user ID from context
func GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: gin.ErrorTypePublic, Meta: "User ID not found in context"}
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, gin.Error{Err: gin.ErrorTypePublic, Meta: "Invalid user ID type"}
}

// GetCurrentUserRole gets the current user role from context
func GetCurrentUserRole(c *gin.Context) (string, error) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", gin.Error{Err: gin.ErrorTypePublic, Meta: "User role not found in context"}
	}

	if roleStr, ok := role.(string); ok {
		return roleStr, nil
	}

	return "", gin.Error{Err: gin.ErrorTypePublic, Meta: "Invalid user role type"}
}