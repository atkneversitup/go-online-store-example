package database

import (
	"github.com/jinzhu/gorm"
)

// Model interface represents the database model
type OrderProduct struct {
	gorm.Model
	// ID        uint `gorm:"primary_key;auto_increment"`
	OrderID   uint
	ProductID uint
	Quantity  int

	// Order   Order
	Product Product
}

// Users is a slice of User
type OrderProducts []OrderProduct
