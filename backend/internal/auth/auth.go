package auth

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ✅ JWT Claims Struct
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ✅ Load JWT Secrets Securely
func getSecret(envVar string, defaultValue string) []byte {
	secret := os.Getenv(envVar)
	if secret == "" {
		log.Printf("⚠️ WARNING: %s is not set, using default value!", envVar)
		secret = defaultValue // Set a default secret (only for development)
	}
	return []byte(secret)
}

var jwtSecret = getSecret("JWT_SECRET", "default-secret-key-should-be-longer-than-this")
var jwtRefreshSecret = getSecret("JWT_REFRESH_SECRET", "default-refresh-key-should-be-longer")

// ✅ Regex for password, username, and email validation
var (
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d@$!%*?&.]{8,64}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,32}$`) // Allows letters, numbers, _ and . (3-32 chars)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ✅ Validate Password Strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 64 {
		return errors.New("password must not exceed 64 characters")
	}
	if !passwordRegex.MatchString(password) {
		return errors.New("password contains invalid characters")
	}

	hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false
	specialChars := "@$!%*?&."

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least 1 uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least 1 lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least 1 number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least 1 special character (@$!%*?&)")
	}

	return nil
}

// ✅ Validate Username
func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return errors.New("username must be 3-32 characters and only contain letters, numbers, underscores, or dots")
	}
	return nil
}

// ✅ Validate Email
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ✅ Hash Password
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// ✅ Check Password Hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("❌ Password verification failed:", err)
		return false
	}
	return true
}

// ✅ Generate JWT Access Token (1 hour expiry)
func GenerateAccessToken(userID uuid.UUID) (string, error) {
	return generateToken(userID.String(), jwtSecret, time.Hour)
}

// ✅ Generate JWT Refresh Token (7 days expiry)
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	return generateToken(userID.String(), jwtRefreshSecret, 7*24*time.Hour)
}

// ✅ Core JWT Token Generation Function
func generateToken(userID string, secret []byte, expiry time.Duration) (string, error) {
	if len(secret) < 32 {
		log.Println("⚠️ WARNING: JWT secret is too short. Use at least 32 characters!")
	}

	expirationTime := time.Now().Add(expiry)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Println("❌ Error signing JWT:", err)
		return "", err
	}

	return signedToken, nil
}

// ✅ Validate JWT Token
func ValidateToken(tokenString string, isRefresh bool) (*Claims, error) {
	secret := jwtSecret
	if isRefresh {
		secret = jwtRefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// ✅ Check if token has expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
