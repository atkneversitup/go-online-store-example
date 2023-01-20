package database

import (
	"github.com/jinzhu/gorm"
)

// Model interface represents the database model
type Order struct {
	gorm.Model
	UserID   int
	StatusID uint

	User          User
	Status        Status
	OrderProducts []OrderProduct `gorm:"foreignkey:OrderID"`
}

// Users is a slice of User
type Orders []Order
