package middleware

import (
	"fmiis/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware interface {
	Handle() gin.HandlerFunc
	RequireRole(allowedRoles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	authService auth.AuthService
}

func NewAuthMiddleware(authService auth.AuthService) AuthMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

func (m *authMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only check cookie for token
		tokenString, err := c.Cookie("jwt_token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"details": "Please log in to access this resource",
			})
			return
		}

		token, err := m.authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is invalid or expired",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		userId, ok := claims["user_id"].(string)
		if !ok || userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Undefined user role",
			})
			return
		}

		c.Set("userID", userId)
		c.Set("role", role)
		c.Next()
	}
}

func (m *authMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
			return
		}

		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: insufficient permissions",
			"details": "Required roles: " + strings.Join(allowedRoles, ", ") +
				", Your role: " + roleStr,
		})
	}
}
