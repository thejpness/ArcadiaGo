package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thejpness/ArcadiaGo/internal/auth"
)

// ✅ AuthMiddleware - Protects routes by requiring authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			log.Println("⚠️ No authentication token found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(token, false)
		if err != nil {
			log.Println("❌ Invalid authentication token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ✅ Store user_id in context instead of email
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
