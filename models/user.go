package models

import (
	"fmt"

	"github.com/msal4/toastnotes/db"
	"gorm.io/gorm"
)

// User is the model representing regular users.
//
// TODO: add auth keys and other stuff relating to authorization and authentication.
type User struct {
	gorm.Model
	ID       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}

// SignUpForm is used to register a new user.
type SignUpForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserSignUp creates a new user record using a SignUpForm.
func UserSignUp(data *SignUpForm) (*User, error) {
	user := User{Name: data.Name, Email: data.Email, Password: data.Password}

	if err := db.DB.Create(&user).Error; err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("could not create user")
	}

	return &user, nil
}
