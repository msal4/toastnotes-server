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
	// init
	godotenv.Load()
	db.Init()
	validation.UseJSONFieldNames()

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintln("migrating failed with error:", err))
	}

	// config
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// router
	router := gin.Default()

	router.Use(middleware.CORS())

	// controllers
	userController := controllers.NewUserController()

	v1 := router.Group("api/v1")
	{
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		v1.POST("/refresh", userController.RefreshTokens)

		authenticated := v1.Group("/", middleware.JWTAuth())
		{
			authenticated.GET("/me", userController.Me)
			authenticated.POST("/change_password", userController.ChangePassword)
		}
	}

	router.Run()
}
