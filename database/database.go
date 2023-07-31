package database

import (
	"fmt"
	"log"
	"os"
	"salamander-smtp/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeThunderDome() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	db := os.Getenv("DB_ID")
	secret := os.Getenv("DB_SECRET")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, db, secret)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(&models.VerifyUserTemplate{})

	return err
}

func MigrateVerifiactionEmailEvent() {
	// gormDB := database.FetchDB()
	DB.AutoMigrate(&models.VerificationEmailEvent{})
}

func FetchDB() *gorm.DB {
	return DB
}
