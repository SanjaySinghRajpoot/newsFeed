package main

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	router "github.com/SanjaySinghRajpoot/newsFeed/router"
	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to News Feed"})
}

func main() {

	// Connect to the database
	config.Connect()

	// Gin router
	r := gin.Default()

	// Home Page endpoint
	r.GET("/", HomepageHandler)

	// All the Routes
	router.GetRoute(r)

	r.Run(":8081")
}
