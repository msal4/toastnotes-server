package models

import (
	"errors"
	"fmt"

	"github.com/msal4/toastnotes/auth"
	"gorm.io/gorm"
)

// User is the model representing standard users.
type User struct {
	Model
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	TokenVersion int    `json:"-" gorm:"default:0"`
	Notes        []Note `json:"notes"`
}

// UserRepository holds all the database operations related to the user.
type UserRepository struct {
	*Repository
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Repository: &Repository{DB: db}}
}

// RegisterUser creates a new user record using a SignUpForm.
func (rep *UserRepository) RegisterUser(data auth.RegisterForm) (*User, error) {
	password, err := auth.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := User{Name: data.Name, Email: data.Email, Password: password}
	if err := rep.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("Failed to create user")
	}
	return &user, nil
}

// RetrieveUser finds the user with the given id.
func (rep *UserRepository) RetrieveUser(id string) (*User, error) {
	var user User
	if err := rep.FindByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

// EmailTaken check if a user has already registered with the given email.
func (rep *UserRepository) EmailTaken(email string) bool {
	err := rep.DB.First(&User{}, "email = ?", email).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}
