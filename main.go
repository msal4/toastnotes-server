package main

import (
	"fmt"

	"github.com/msal4/toastnotes/db"
	"github.com/msal4/toastnotes/models"
)

func main() {
	db.Init()

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Sprintln("migrating User failed with error:", err))
	}

	u, err := models.UserSignUp(&models.SignUpForm{
		Name:     "Mohammed Salman",
		Email:    "msal4@outlook.com",
		Password: "password",
	})

	if err != nil {
		return
	}

	fmt.Printf("%+v\n", u)
}
