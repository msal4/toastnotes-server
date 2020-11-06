package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
func OpenConnection(dsn string, lgr logger.Interface) (*gorm.DB, error) {
	if lgr == nil {
		lgr = logger.Default
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: lgr,
	})
	if err != nil {
		return nil, err
	}

	// Create the uuid extension to generate uuids for the id field in models.
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, errors.New("Could not create extension \"uuid-ossp\"")
	}

	if err := db.AutoMigrate(&User{}, &Note{}); err != nil {
		return nil, err
	}

	return db, nil
}

// Paginate paginates the given request context using scopes.
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
		switch {
		case pageSize > settings.MaxPageSize:
			pageSize = settings.MaxPageSize
		case pageSize <= 0:
			pageSize = settings.PageSize
		}

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}
