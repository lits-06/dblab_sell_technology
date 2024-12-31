package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lits-06/sell_technology/internal/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect () {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Info,
			Colorful: true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	DB = db
	log.Println("Connected to database!")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v\n", err)
	}
}