package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the database instance.
var DB *gorm.DB

// Init connects to the database.
func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open("postgres://msal@localhost:5432/toast"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Create the uuid extension to generate uuids for the id field in models.
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		fmt.Println("Could not create extension \"uuid-ossp\"")
	}
}
