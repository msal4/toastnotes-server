package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/controllers"
	"github.com/msal4/toastnotes/db"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// init
	godotenv.Load()
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	validation.UseJSONFieldNames()

	if err := conn.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintln("migrating failed with error:", err))
	}

	// config
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// router
	router := controllers.SetupRouter(conn)
	if err := router.Run(); err != nil {
		panic(err)
	}
}
