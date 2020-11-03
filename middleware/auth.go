package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/auth"
)

// JWTAuth is the auth middleware that handles jwt authentication.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		abortUnauthorized := func() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		tokenStr, err := c.Cookie(auth.AccessTokenKey)
		if err != nil {
			abortUnauthorized()
			return
		}

		claims := auth.AccessTokenClaims{}
		token, err := auth.ParseToken(tokenStr, &claims)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				abortUnauthorized()
				return
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		if !token.Valid {
			abortUnauthorized()
			return
		}

		c.Set(auth.UserIDKey, claims.UserID)

		c.Next()
	}
}
