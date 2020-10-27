package models

import (
	"errors"
	"fmt"

	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/db"
	"gorm.io/gorm"
)

// User is the model representing standard users.
type User struct {
	Base
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	TokenVersion int    `json:"-" gorm:"default:0"`
}

// RegisterUser creates a new user record using a SignUpForm.
func (u *User) RegisterUser(data auth.RegisterForm) (*User, error) {
	password, err := auth.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := User{Name: data.Name, Email: data.Email, Password: password}
	if err := db.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("Failed to create user")
	}
	return &user, nil
}

// RetrieveUser finds the user with the given id.
func (u *User) RetrieveUser(id string) (*User, error) {
	var user User
	if err := FindByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

// EmailTaken check if a user has already registered with the given email.
func (u *User) EmailTaken(email string) bool {
	err := db.DB.First(&User{}, "email = ?", email).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}
