package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real notification"})
}

func main() {

	// Gin router
	r := gin.Default()

	// Home Page endpoint
	r.GET("/", HomepageHandler)

	r.Run(":8081")
}
