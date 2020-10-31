package main

import (
	"os"

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
	db, err := models.OpenConnection(os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
	}

	validation.UseJSONFieldNames()

	// config
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// router
	router := controllers.SetupRouter(db)
	if err := router.Run(); err != nil {
		panic(err)
	}
}
