package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejpness/ArcadiaGo/internal/database"
	"github.com/thejpness/ArcadiaGo/internal/models"
)

// RequestEmailVerification handles email change requests by sending a confirmation email
func RequestEmailVerification(c *gin.Context) {
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
		NewEmail string `json:"new_email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("❌ Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Println("🔍 Requested email change for:", userID, "to:", req.NewEmail)

	// Ensure the new email is not already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.NewEmail).First(&existingUser).Error; err == nil {
		log.Println("❌ Email already registered:", req.NewEmail)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Generate a verification token
	token := uuid.New().String()

	// Store the email verification request
	emailChange := models.UserEmailChange{
		ID:        uuid.New(),
		UserID:    userID,
		NewEmail:  req.NewEmail,
		Token:     token,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&emailChange).Error; err != nil {
		log.Println("❌ Failed to create email verification request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email verification request"})
		return
	}

	// Send confirmation email via MailHog
	err = SendEmail(req.NewEmail, "Confirm Email Change",
		fmt.Sprintf("Click here to confirm your email change: http://localhost:8080/confirm-email?token=%s", token))
	if err != nil {
		log.Println("❌ Failed to send email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	log.Println("✅ Email change request stored and verification email sent")
	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// ConfirmEmailVerification verifies the token and updates the user's email
func ConfirmEmailVerification(c *gin.Context) {
	token := c.Query("token")

	// Check if the token exists in the user_email_changes table
	var request models.UserEmailChange
	if err := database.DB.Where("token = ?", token).First(&request).Error; err != nil {
		log.Println("❌ Invalid or expired token:", token)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token"})
		return
	}

	log.Println("🔍 Found email change request for user:", request.UserID, "New Email:", request.NewEmail)

	// Start a transaction to ensure atomicity
	tx := database.DB.Begin()

	// Fetch the user by ID
	var user models.User
	if err := tx.Where("id = ?", request.UserID).First(&user).Error; err != nil {
		log.Println("❌ User not found in users table:", request.UserID, err)
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	log.Println("👀 Current Email before update:", user.Email)

	// Update the user's email
	if err := tx.Model(&models.User{}).Where("id = ?", request.UserID).Update("email", request.NewEmail).Error; err != nil {
		log.Println("❌ Failed to update email for user:", request.UserID, err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
		return
	}

	// Fetch the user after the update to confirm the change
	var updatedUser models.User
	if err := tx.Where("id = ?", request.UserID).First(&updatedUser).Error; err == nil {
		log.Println("✅ Updated Email:", updatedUser.Email)
	} else {
		log.Println("❌ Failed to fetch updated user email:", err)
	}

	// Delete the email change request only if the update was successful
	if err := tx.Where("token = ?", token).Delete(&models.UserEmailChange{}).Error; err != nil {
		log.Println("❌ Failed to delete email verification request:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize email update"})
		return
	}

	// Commit the transaction
	tx.Commit()

	log.Println("✅ Email updated successfully for user:", request.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}

// SendEmail sends an email using MailHog (SMTP)
func SendEmail(to, subject, body string) error {
	smtpHost := "localhost"
	smtpPort := "1025" // MailHog SMTP port
	from := "no-reply@arcadiago.dev"

	// Construct email headers
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	// Send email using MailHog
	err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, []string{to}, msg)
	if err != nil {
		log.Println("❌ Failed to send email:", err)
		return err
	}

	log.Println("📧 Email sent successfully to:", to)
	return nil
}
