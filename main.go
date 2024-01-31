package main

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	router "github.com/SanjaySinghRajpoot/newsFeed/router"
	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real notification"})
}

func main() {

	// Connect to the database
	config.Connect()

	// Gin router
	r := gin.Default()

	router.GetRoute(r)

	// Home Page endpoint
	r.GET("/", HomepageHandler)

	r.Run(":8081")
}
