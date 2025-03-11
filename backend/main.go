package main

import (
	"log"
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
	// ✅ Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: No .env file found")
	}

	// ✅ Initialize the database
	database.InitDB()
	if database.DB == nil {
		log.Fatal("❌ Failed to connect to the database")
	}

	// ✅ Create a new Gin router
	r := gin.New()

	// ✅ Middleware Stack
	r.Use(gin.Logger())            // Logs all requests
	r.Use(gin.Recovery())          // Prevents crashes from panics
	r.Use(middleware.CORSConfig()) // Enables CORS
	r.Use(setupSecurityHeaders())  // Adds security headers
	r.Use(setupRateLimiter())      // Enables Rate Limiting

	// ✅ Default Root Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is running"})
	})

	// ✅ Define Public API Routes
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginUser)
	r.POST("/logout", handlers.LogoutUser)

	// ✅ Protected Routes (Require Authentication)
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware()) // ✅ Secure all endpoints below
	{
		authenticated.GET("/user", handlers.GetUserProfile)   // ✅ Fetch user profile
		authenticated.POST("/refresh", handlers.RefreshToken) // ✅ Refresh Access Token
	}

	// ✅ Start the server
	log.Println("🚀 Server running on port 8080")
	log.Fatal(r.Run(":8080"))
}

// ✅ Correct Rate Limiting (5 requests per minute per IP)
func setupRateLimiter() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute, // ✅ Set rate limit expiration
	})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}) // ✅ Track request IP

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
