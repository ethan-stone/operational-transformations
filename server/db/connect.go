package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	database.AutoMigrate(new(Document))
	database.AutoMigrate(new(Operation))

	DB = database
}
