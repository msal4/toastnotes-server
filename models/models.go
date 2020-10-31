package models

import (
	"time"

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
