package config

import (
	"log"
	"os"

	"example-project.com/event-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URI")

	if dsn == "" {
		log.Fatal("Env variable belum diisi")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed Connect to Database", err)
	}

	err = database.AutoMigrate(&models.Event{}, models.User{})
	if err != nil {
		log.Fatal("Database Migration Failed", err)
	}

	DB = database
	log.Println("Success Connect to Database")

}
