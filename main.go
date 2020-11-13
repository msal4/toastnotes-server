package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/controllers"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// init
	godotenv.Load()
	databaseURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := models.OpenConnection(databaseURI, nil)
	if err != nil {
		panic(err)
	}

	validation.UseJSONFieldNames()

	// config
	auth.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// router
	router := controllers.SetupRouter(db)

	// middlewares
	router.Use(gin.Logger())

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("cert-cache"),
		HostPolicy: autocert.HostWhitelist("api.toast.msal.dev"),
		Email:      "msal4@outlook.com",
	}

	server := http.Server{
		Addr:    ":443",
		Handler: router,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}
