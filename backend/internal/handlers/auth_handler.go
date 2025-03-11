package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thejpness/ArcadiaGo/internal/auth"
	"github.com/thejpness/ArcadiaGo/internal/database"
	"github.com/thejpness/ArcadiaGo/internal/models"
)

// ✅ Register a new user with password validation
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ Check if user already exists
	var existingUser models.User
	err := database.DB.Get(&existingUser, "SELECT * FROM users WHERE email=$1", user.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// ✅ Validate password before hashing
	if err := auth.ValidatePassword(user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Hash password before storing
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	user.Password = hashedPassword

	_, err = database.DB.NamedExec(`INSERT INTO users (email, password) VALUES (:email, :password)`, user)
	if err != nil {
		log.Println("❌ Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// ✅ Login user, set JWT cookies, and issue a refresh token
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var dbUser models.User
	err := database.DB.Get(&dbUser, "SELECT * FROM users WHERE email=$1", user.Email)
	if err != nil || auth.CheckPassword(dbUser.Password, user.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// ✅ Generate JWT Access & Refresh Tokens
	accessToken, err := auth.GenerateAccessToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token generation failed"})
		return
	}

	// ✅ Set Secure HttpOnly Cookies
	c.SetCookie("auth_token", accessToken, 3600, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, int(7*24*time.Hour.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ✅ Logout user by clearing authentication and refresh token cookies
func LogoutUser(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ✅ Refresh Access Token using a valid Refresh Token
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token found"})
		return
	}

	// ✅ Corrected: Capture both `claims` and `err` properly
	claims, err := auth.ValidateToken(refreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// ✅ Generate a new Access Token
	newAccessToken, err := auth.GenerateAccessToken(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return
	}

	// ✅ Set new Secure HttpOnly Access Token Cookie
	c.SetCookie("auth_token", newAccessToken, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// ✅ Get user profile based on authenticated JWT
func GetUserProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": email})
}
