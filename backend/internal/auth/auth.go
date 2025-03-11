package auth

import (
	"errors"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// ✅ Struct to store JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// ✅ Securely load JWT Secret & Signing Method from environment variables
var jwtSecret []byte
var jwtRefreshSecret []byte
var jwtSigningMethod *jwt.SigningMethodHMAC
var accessTokenExpiry time.Duration
var refreshTokenExpiry time.Duration

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: No .env file found, using default JWT secret.")
	}

	secret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	if secret == "" || refreshSecret == "" {
		log.Fatal("❌ JWT_SECRET or JWT_REFRESH_SECRET is missing! Set them in .env")
	}

	jwtSecret = []byte(secret)
	jwtRefreshSecret = []byte(refreshSecret)

	// ✅ Allow configurable signing method (default: HS256)
	signingMethod := os.Getenv("JWT_SIGNING_METHOD")
	if signingMethod == "HS512" {
		jwtSigningMethod = jwt.SigningMethodHS512
	} else {
		jwtSigningMethod = jwt.SigningMethodHS256
	}

	// ✅ Set configurable expiration times (default: 1hr access, 7d refresh)
	accessTokenExpiry = time.Hour
	refreshTokenExpiry = 7 * 24 * time.Hour
}

// ✅ Secure Password Validation Function (Fixed for Go)
func ValidatePassword(password string) error {
	// ✅ Manually enforce rules instead of using unsupported lookaheads
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least 1 uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least 1 lowercase letter")
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return errors.New("password must contain at least 1 number")
	}
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
		return errors.New("password must contain at least 1 special character (@$!%*?&)")
	}

	return nil
}

// ✅ Hash password before storing in DB
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err // ❌ Reject weak passwords
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// ✅ Compare stored hash with user input password
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("❌ Password validation failed:", err)
	}
	return err
}

// ✅ Generate JWT access token
func GenerateAccessToken(email string) (string, error) {
	return generateToken(email, jwtSecret, accessTokenExpiry)
}

// ✅ Generate JWT refresh token
func GenerateRefreshToken(email string) (string, error) {
	return generateToken(email, jwtRefreshSecret, refreshTokenExpiry)
}

// ✅ Core function for generating tokens
func generateToken(email string, secret []byte, expiry time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiry)

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwtSigningMethod, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Println("❌ Error signing JWT:", err)
		return "", err
	}

	return signedToken, nil
}

// ✅ Validate JWT token (supports both access & refresh)
func ValidateToken(tokenString string, isRefresh bool) (*Claims, error) {
	secret := jwtSecret
	if isRefresh {
		secret = jwtRefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("❌ Invalid JWT signing method:", token.Header["alg"])
			return nil, errors.New("invalid signing method")
		}
		return secret, nil
	})

	if err != nil {
		log.Println("❌ Token validation failed:", err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("❌ Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
