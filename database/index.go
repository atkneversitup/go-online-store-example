package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	var err error
	log.Println("Connecting to database")
	DB, err = gorm.Open("postgres", "user=postgres dbname=go-test3 password=123456 sslmode=disable")
	if err != nil {
		panic("failed to connect to database")
	}
	log.Println("Database successfully connected")
	// Automatically migrate the schema
	log.Println("Configuring database")
	DB.AutoMigrate(User{})
	DB.AutoMigrate(Status{})
	DB.AutoMigrate(Product{})
	DB.AutoMigrate(Order{})
	DB.AutoMigrate(OrderProduct{})
	log.Println("Database successfully configured")
}
