package models

import (
	"time"

	"github.com/msal4/toastnotes/db"
	"gorm.io/gorm"
)

// Model is the base model.
type Model struct {
	ID        string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"index"`
}

// FindByID finds a record with the given id.
func FindByID(v interface{}, id string) error {
	return db.DB.First(v, "id = ?", id).Error
}
