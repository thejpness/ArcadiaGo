package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejpness/ArcadiaGo/internal/auth"
	"github.com/thejpness/ArcadiaGo/internal/database"
	"github.com/thejpness/ArcadiaGo/internal/models"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Println("🔍 Received password:", user.Password) // ✅ Debug received password

	// ✅ Ensure Email Uniqueness
	var existingUser models.User
	err := database.DB.Get(&existingUser, "SELECT id FROM users WHERE email=$1", user.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// ✅ Validate Password Strength
	if err := auth.ValidatePassword(user.Password); err != nil {
		log.Println("❌ Password validation failed:", err) // ✅ Debug validation
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Hash Password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	user.ID = uuid.New()
	user.Password = hashedPassword

	// ✅ Insert User into Database
	_, err = database.DB.NamedExec(
		`INSERT INTO users (id, email, password, created_at) VALUES (:id, :email, :password, NOW())`,
		user,
	)
	if err != nil {
		log.Println("❌ Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// ✅ Login user and issue tokens
func LoginUser(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ Fetch User by Email
	var user models.User
	err := database.DB.Get(&user, "SELECT id, password FROM users WHERE email=$1", request.Email)
	if err != nil || !auth.CheckPassword(user.Password, request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// ✅ Generate JWT Access & Refresh Tokens using UUID
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh token generation failed"})
		return
	}

	// ✅ Set Secure HttpOnly Cookies
	c.SetCookie("auth_token", accessToken, 3600, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, int(7*24*time.Hour.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ✅ Logout user by clearing authentication & refresh token cookies
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

	// ✅ Validate the Refresh Token
	claims, err := auth.ValidateToken(refreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// ✅ Convert UserID back to UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// ✅ Generate a new Access Token
	newAccessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token"})
		return
	}

	// ✅ Set new Secure HttpOnly Access Token Cookie
	c.SetCookie("auth_token", newAccessToken, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// ✅ Fetch User Profile using UUID stored in JWT
func GetUserProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var user struct {
		Email string `json:"email"`
	}
	err = database.DB.Get(&user, "SELECT email FROM users WHERE id=$1", userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
