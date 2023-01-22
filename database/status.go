package database

import (
	"github.com/jinzhu/gorm"
)

// User struct represents the user model
type Status struct {
	// ID int
	gorm.Model
	Name string
}
