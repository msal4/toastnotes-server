package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Credentials are the needed credentials to log a user in.
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterForm is used to register a new user.
type RegisterForm struct {
	Credentials
	Name string `json:"name" binding:"required"`
}

// ChangePasswordForm is intended to be used when a user requests to change their password.
type ChangePasswordForm struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=8"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// AccessTokenClaims ...
type AccessTokenClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

// RefreshTokenClaims ...
type RefreshTokenClaims struct {
	UserID       string `json:"userId"`
	TokenVersion int    `json:"tokenVersion"`
	jwt.StandardClaims
}

const (
	// AccessTokenAge is the token age in seconds.
	AccessTokenAge = 300 // = 5 minutes

	// RefreshTokenAge is the refresh token age in seconds.
	RefreshTokenAge = 2.628e6 // = 1 month

	// UserIDKey is the key used to set the user id in gin context.
	UserIDKey = "userId"

	// RefreshTokenKey is the key used to set the refresh token cookie
	RefreshTokenKey = "refresh_token"

	// AccessTokenKey is the key used to set the refresh token cookie
	AccessTokenKey = "token"

	// PasswordHashCost is the cost used for hashing the user password.
	PasswordHashCost = 11
)

// JWTSecret is the secret jwt key used to create tokens.
var JWTSecret = os.Getenv("JWT_SECRET")

// GenerateAccessToken generates an access token using the refresh token.
func GenerateAccessToken(userID string) (string, error) {
	claims := &AccessTokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenAge * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates an access token using the refresh token.
func GenerateRefreshToken(userID string, version int) (string, error) {
	claims := &RefreshTokenClaims{
		UserID:       userID,
		TokenVersion: version,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenAge * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// HashPassword hashes the password string.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
	return string(hash), err
}

// PasswordMatch checks if the password matches the hash.
func PasswordMatch(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
