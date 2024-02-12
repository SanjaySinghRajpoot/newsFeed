package config

import (
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

// use godot package to load/read the .env file and
// return the value of the key
// func goDotEnvVariable(key string) string {

// 	// load .env file
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	return os.Getenv(key)
// }

func Connect() {

	// dsn := os.Getenv("DATABASE_URL")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// err = db.Migrator().DropTable(models.User{}, models.Post{}, &models.Comment{}, &models.Follower{})
	// if err != nil {
	// 	log.Fatal("Table dropping failed")
	// }

	fmt.Println(dsn)

	// err1 := db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Follower{})

	// if err1 != nil {
	// 	log.Fatal("Migration Error ", err1)
	// }

	DB = db
}
