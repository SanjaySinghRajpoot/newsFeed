package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//   dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
//   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

var DB *gorm.DB

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func Connect() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	adminEmail := os.Getenv("ADMIN_EMAIL")

	fmt.Print(adminEmail)

	// Check if admin already exist
	row := db.Where("email = ?", adminEmail).Row()

	if row == nil {
		// Create a defualt admin User
		user := models.User{
			Name:     os.Getenv("ADMIN_NAME"),
			Email:    os.Getenv("ADMIN_EMAIL"),
			Password: os.Getenv("ADMIN_PASSWORD"),
			Admin:    true,
		}
		db.Create(&user)
	}

	DB = db
}
