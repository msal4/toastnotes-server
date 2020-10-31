package models

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Model is the base model.
type Model struct {
	ID        string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"index"`
}

// Repository is the base repository.
type Repository struct {
	DB *gorm.DB
}

// FindByID finds a record with the given id.
func (rep *Repository) FindByID(v interface{}, id string) error {
	return rep.DB.First(v, "id = ?", id).Error
}

// OpenConnection opens a db connections using the provided uri.
func OpenConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Create the uuid extension to generate uuids for the id field in models.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, errors.New("Could not create extension \"uuid-ossp\"")
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, nil
}
