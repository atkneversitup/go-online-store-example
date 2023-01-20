package database

import (
	"github.com/jinzhu/gorm"
)

// User struct represents the user model
type User struct {
	gorm.Model
	Name     string
	Username string
	Password string
	Orders   []Order
}

// Users is a slice of User
type Users []User
