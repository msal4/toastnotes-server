package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//CORS (Cross-Origin Resource Sharing).
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	origins := os.Getenv("ALLOW_ORIGINS")
	if origins != "" {
		config.AllowOrigins = strings.Split(origins, ",")
		fmt.Println("origins:", config.AllowOrigins)
	}
	config.AllowCredentials = true

	return cors.New(config)
}
