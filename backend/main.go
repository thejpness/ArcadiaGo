package main

import (
	"log"
	"os"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thejpness/ArcadiaGo/internal/database"
	"github.com/thejpness/ArcadiaGo/internal/handlers"
	"github.com/thejpness/ArcadiaGo/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: No .env file found, using system environment variables")
	}

	// Initialize and Migrate Database
	database.InitDB()
	database.Migrate()

	if database.DB == nil {
		log.Fatal("❌ Failed to connect to the database")
	}

	// Get server port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new Gin router
	r := gin.New()

	// Middleware Stack
	r.Use(gin.Logger())            // Logs all requests
	r.Use(gin.Recovery())          // Prevents crashes from panics
	r.Use(middleware.CORSConfig()) // Enables CORS
	r.Use(setupSecurityHeaders())  // Adds security headers
	r.Use(setupRateLimiter())      // Enables Rate Limiting

	// Default Root Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is running"})
	})

	// ✅ Public API Routes (No Authentication Required)
	publicRoutes := r.Group("/")
	{
		publicRoutes.POST("/register", handlers.RegisterUser)
		publicRoutes.POST("/login", handlers.LoginUser)
		publicRoutes.POST("/logout", handlers.LogoutUser)
		publicRoutes.GET("/confirm-email", handlers.ConfirmEmailVerification) // Fixed function name
	} // ✅ Closing bracket was missing

	// ✅ Protected Routes (Require Authentication)
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware()) // Secure all endpoints below
	{
		// User Profile
		authenticated.GET("/user", handlers.GetUserProfile)
		authenticated.POST("/refresh", handlers.RefreshToken)

		// User Management
		authenticated.POST("/update-email", handlers.RequestEmailChange) // Request email change
		authenticated.POST("/update-password", handlers.UpdatePassword)  // Change password
		authenticated.POST("/update-username", handlers.UpdateUsername)  // Change username

		// Account Management
		authenticated.POST("/delete-account", handlers.SoftDeleteUser) // Soft delete account
		authenticated.POST("/restore-account", handlers.RestoreUser)   // Restore deleted account

		// Session Management
		authenticated.GET("/active-sessions", handlers.GetActiveSessions) // List active sessions
		authenticated.POST("/logout-session", handlers.LogoutSession)     // Logout from a specific session
	}

	// Start the server
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// ✅ Rate Limiting (10 requests per minute per IP)
func setupRateLimiter() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute, // Set rate limit expiration
	})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}) // Track request IP

	return tollbooth_gin.LimitHandler(lmt)
}

// ✅ Security Headers Middleware
func setupSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
