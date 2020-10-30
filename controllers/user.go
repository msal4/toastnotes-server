package controllers

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/db"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/utils"
	"github.com/msal4/toastnotes/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register a new user.
func Register(c *gin.Context) {
	// Validate the form
	var form auth.RegisterForm
	if errs := shouldBindJSON(c, &form); errs != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, *errs)
		return
	}

	// Check if the email is taken.
	if models.EmailTaken(form.Email) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Err("A user with this email already exists"))
		return
	}

	// Create the user.
	user, err := models.RegisterUser(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err("Failed to register user"))
		return
	}

	generateTokens(c, user.ID, user.TokenVersion, user)
}

// Login a user.
func Login(c *gin.Context) {
	var credentials auth.Credentials
	if errs := shouldBindJSON(c, &credentials); errs != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, *errs)
		return
	}

	var user models.User
	if err := db.DB.First(&user, "email = ?", credentials.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("User not found"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Failed to find the user"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Wrong email or password"))
		return
	}

	generateTokens(c, user.ID, user.TokenVersion, gin.H{"message": "Login successful"})
}

// ChangePassword takes the current password for the authenticated user and allows them to set a new
// password.
func ChangePassword(c *gin.Context) {
	var form auth.ChangePasswordForm
	if err := shouldBindJSON(c, &form); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}

	user, err := models.RetrieveUser(c.GetString(auth.UserIDKey))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("User not found"))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Failed to find the user"))
		return
	}

	if !auth.PasswordMatch(user.Password, form.CurrentPassword) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Wrong password"))
		return
	}

	if form.CurrentPassword == form.NewPassword {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.Err("Please use a different password"))
		return
	}

	hash, err := auth.HashPassword(form.NewPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Failed to update password"))
		return
	}

	if err := db.DB.Model(&user).Update("password", hash).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Failed to update password"))
		return
	}

	c.JSON(http.StatusOK, utils.Msg("Password updated"))
}

// Me retrieves the authenticated user details.
func Me(c *gin.Context) {
	userID := c.GetString(auth.UserIDKey)

	user, err := models.RetrieveUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("User not found"))
			return
		}

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, user)
}

// RefreshToken uses the refresh token to generate an access token and regenerates refresh_token.
func RefreshToken(c *gin.Context) {
	tokenStr, err := c.Cookie(auth.RefreshTokenKey)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Unauthorized"))
		return
	}

	claims := auth.RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return auth.JWTSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Unauthorized"))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Err("Bad request"))
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Unauthorized"))
		return
	}

	user, err := models.RetrieveUser(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("User not found"))
			return
		}

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if user.TokenVersion != claims.TokenVersion {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Err("Token version mismatch"))
		return
	}

	generateTokens(c, user.ID, user.TokenVersion, utils.Msg("Tokens refreshed"))
}

func shouldBindJSON(c *gin.Context, obj interface{}) *gin.H {
	if err := c.ShouldBindJSON(obj); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			return &gin.H{"errors": validation.DescriptiveErrors(verr)}
		}
		return &gin.H{"error": "Invalid form"}
	}
	return nil
}

func generateTokens(c *gin.Context, userID string, tokenVersion int, resp interface{}) {
	tokenStr, err := auth.GenerateAccessToken(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	refreshTokenStr, err := auth.GenerateRefreshToken(userID, tokenVersion)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SetCookie(auth.AccessTokenKey, tokenStr, auth.AccessTokenAge, "/", c.Request.URL.Hostname(), false, true)
	c.SetCookie(auth.RefreshTokenKey, refreshTokenStr, auth.RefreshTokenAge, "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, resp)
}