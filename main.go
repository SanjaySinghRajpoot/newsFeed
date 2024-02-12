package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	router "github.com/SanjaySinghRajpoot/newsFeed/router"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
	sentimentanalysis "github.com/SanjaySinghRajpoot/newsFeed/utils/sentimentAnalysis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/robfig/cron.v2"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to News Feed"})
}

func CRONjobs() {
	// start the cron job
	cronJob := cron.New()

	// Cron Job that will check for state=false for Notifications
	cronJob.AddFunc("@every 5s", func() {

		fmt.Println("cron working")

		sentimentanalysis.UpdateSentimentAnalysisCRON()
	})

	cronJob.Start()
}

func main() {

	// Connect to the database
	config.Connect()

	// start the CRON JOB
	CRONjobs()

	// Gin router
	r := gin.Default()

	// Home Page endpoint
	r.GET("/", HomepageHandler)

	password := EnvVariable("PASSWORD")

	// Redis Cache Setup
	redis.RedisClient = redis.SetUpRedis(password)

	// All the Routes
	router.GetRoute(r)

	r.Run(":8081")
}

// return the value of the key
func EnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
