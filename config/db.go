package config

import (
	"log"
	"os"

	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	// dsn := "host=localhost user=postgres password=postgres dbname=newsfeed port=5432 sslmode=disable"

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// err = db.Migrator().DropTable(models.User{}, models.Post{}, models.UserTest{})
	// if err != nil {
	// 	log.Fatal("Table dropping failed")
	// }

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Follower{})

	if err != nil {
		log.Fatal("Migration Error", err.Error())
	}

	DB = db
}
