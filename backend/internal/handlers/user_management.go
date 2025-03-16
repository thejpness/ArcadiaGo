package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejpness/ArcadiaGo/internal/auth"
	"github.com/thejpness/ArcadiaGo/internal/database"
	"github.com/thejpness/ArcadiaGo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// ✅ Validate and Update Password
func UpdatePassword(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		log.Println("❌ User not found:", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify Old Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		log.Println("❌ Incorrect old password for user:", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
		return
	}

	// ✅ Hash new password using `auth.HashPassword`
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update password
	database.DB.Model(&user).Update("password", hashedPassword)

	log.Println("✅ Password updated successfully for user:", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// ✅ Update Username
func UpdateUsername(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		NewUsername string `json:"new_username"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Println("🔍 Requested new username:", req.NewUsername)

	// ✅ Ensure Username is Unique
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.NewUsername).First(&existingUser).Error; err == nil {
		log.Println("❌ Username already exists:", req.NewUsername)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// ✅ Update username
	err = database.DB.Model(&models.User{}).Where("id = ?", userID).Update("username", req.NewUsername).Error
	if err != nil {
		log.Println("❌ Failed to update username:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	log.Println("✅ Username updated successfully to:", req.NewUsername)
	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

// ✅ Request Email Change
func RequestEmailChange(c *gin.Context) {
	RequestEmailVerification(c) // ✅ Uses function from `email_verification.go`
}

// ✅ Confirm Email Change
func ConfirmEmailChange(c *gin.Context) {
	ConfirmEmailVerification(c) // ✅ Uses function from `email_verification.go`
}

// ✅ Soft Delete Account
func SoftDeleteUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := database.DB.Delete(&models.User{}, userID).Error; err != nil {
		log.Println("❌ Failed to delete account:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	log.Println("✅ Account soft deleted for user:", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted (soft delete)"})
}

// ✅ Restore Account
func RestoreUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := database.DB.Unscoped().Model(&models.User{}).Where("id = ?", userID).Update("deleted_at", nil).Error; err != nil {
		log.Println("❌ Failed to restore account:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore account"})
		return
	}

	log.Println("✅ Account restored successfully for user:", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Account restored successfully"})
}

// ✅ Get Active Sessions
func GetActiveSessions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var sessions []models.UserSession
	if err := database.DB.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		log.Println("❌ Failed to retrieve sessions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
		return
	}

	log.Println("✅ Retrieved active sessions for user:", userID)
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// ✅ Logout from a Specific Session
func LogoutSession(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		log.Println("❌ user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		SessionID uuid.UUID `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := database.DB.Where("user_id = ? AND id = ?", userID, req.SessionID).Delete(&models.UserSession{}).Error; err != nil {
		log.Println("❌ Failed to log out session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out session"})
		return
	}

	log.Println("✅ Session logged out successfully for user:", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Session logged out successfully"})
}
