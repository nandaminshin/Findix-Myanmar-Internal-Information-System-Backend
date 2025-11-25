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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		token, err := m.authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims) // jwt.MapClaims is map[string]interface{}
		if !ok {
			// Try casting to jwt.MapClaims if gin.MapClaims fails or if they are different types in this context
			// But usually they are compatible. Let's be safe and use type assertion on interface{}
			// Re-parsing claims might be safer if we want to be strict, but ValidateToken returns *jwt.Token
			// Let's assume standard jwt.MapClaims
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			// If role is missing, default to user or handle as error.
			// For now, let's assume it might be missing for old tokens, but we want to enforce it.
			role = "glob"
		}

		c.Set("userID", userID)
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

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
	}
}
