package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/controllers"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// init
	godotenv.Load()

	// connect to db
	db, err := models.OpenConnection(os.Getenv("DATABASE_URL"), nil)
	if err != nil {
		panic(err)
	}

	validation.UseJSONFieldNames()

	// config
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// router
	router := controllers.SetupRouter(db)

	// middleware
	router.Use(gin.Logger())

	if err := router.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
