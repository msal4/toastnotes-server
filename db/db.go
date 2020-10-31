package db

import (
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database.
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URI")
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Create the uuid extension to generate uuids for the id field in models.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, errors.New("Could not create extension \"uuid-ossp\"")
	}

	return db, nil
}
