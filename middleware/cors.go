package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS (Cross-Origin Resource Sharing).
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	originsStr := os.Getenv("ALLOW_ORIGINS")

	if originsStr != "" && originsStr != "*" {
		origins := []string{}
		for _, o := range strings.Split(originsStr, ",") {
			origins = append(origins, strings.Trim(o, " "))
		}
		config.AllowOrigins = origins
		config.AllowCredentials = true
	} else {
		config.AllowAllOrigins = true
	}

	return cors.New(config)
}
