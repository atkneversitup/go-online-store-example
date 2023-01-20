package database

import (
	"github.com/jinzhu/gorm"
)

// Product struct represents the product model
type Product struct {
	gorm.Model
	ID    uint `gorm:"primary_key;auto_increment"`
	Name  string
	Price float32

	OrderProducts []OrderProduct
}

// Users is a slice of Product
type Products []Product
