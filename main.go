package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/controllers"
	"github.com/msal4/toastnotes/db"
	"github.com/msal4/toastnotes/middleware"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	godotenv.Load()
	db.Init()

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintln("migrating failed with error:", err))
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// get secret key from env
	if key := os.Getenv("JWT_SECRET"); key != "" {
		auth.JWTSecret = []byte(key)
	}

	validation.UseJSONFieldNames()

	router := gin.Default()

	router.Use(middleware.CORS())

	v1 := router.Group("api/v1")
	{
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)
		v1.POST("/refresh", controllers.RefreshToken)

		authenticated := v1.Group("/", middleware.JWTAuth())
		{
			authenticated.GET("/me", controllers.Me)
			authenticated.POST("/change_password", controllers.ChangePassword)
		}
	}

	router.Run()
}
