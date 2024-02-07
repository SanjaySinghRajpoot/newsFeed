package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/pagination"
	limiter "github.com/SanjaySinghRajpoot/newsFeed/utils/rateLimiter"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/validations.go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CreatePost creates a post
func CreatePost(c *gin.Context) {
	// Get user input from request
	var userInput struct {
		Content string `json:"content" binding:"required,min=2,max=200"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// Create a post
	UserID := helpers.GetAuthUser(c).ID

	// rate limit
	err := limiter.RateLimiter(UserID)

	if err != nil {
		// Return the post
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": err,
		})

		return
	}

	post := models.Post{
		Content: userInput.Content,
		UserID:  UserID,
	}

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		result := config.DB.Create(&post)

		if result.Error != nil {
			formatError.InternalServerError(c, result.Error)
			return
		}
	}()

	go func() {
		defer wg.Done()
		msg, er := redis.SetPostCache(UserID, post)
		if er != nil {
			fmt.Println(msg)
			formatError.InternalServerError(c, er)
			return
		}
	}()

	go func() {
		defer wg.Done()
		err = utils.SendNotification(post)
		if err != nil {
			formatError.InternalServerError(c, err)
			// return
		}
	}()

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"message": post,
	})
}

// GetPosts gets all the post
func GetPosts(c *gin.Context) {
	// Get all the posts
	var posts []models.Post

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, _ := strconv.Atoi(perPageStr)

	preloadFunc := func(query *gorm.DB) *gorm.DB {
		return query.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		})
	}

	result, err := pagination.Paginate(config.DB, page, perPage, preloadFunc, &posts)

	if err != nil {
		formatError.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": result,
	})
}

// ShowPost finds a post by ID
func ShowPost(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	// Find the post
	var post models.Post
	result := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Select("id, post_id, user_id, body, created_at")
	}).First(&post, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// EditPost finds a post by ID
func EditPost(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	// Find the post
	var post models.Post
	result := config.DB.First(&post, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func UpdatePost(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	// Get the data from request body
	var userInput struct {
		Content string `json:"content" binding:"required,min=2,max=200"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find the post by id
	var post models.Post
	result := config.DB.First(&post, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Prepare data to update
	authID := helpers.GetAuthUser(c).ID
	updatePost := models.Post{
		Content: userInput.Content,
		UserID:  authID,
	}

	// Update the post
	result = config.DB.Model(&post).Updates(&updatePost)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// Return the post

	c.JSON(http.StatusOK, gin.H{
		"post": updatePost,
	})
}

func DeletePost(c *gin.Context) {
	// Get the id from the url
	id := c.Param("id")
	var post models.Post

	result := config.DB.First(&post, id)
	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Delete the post
	config.DB.Delete(&post)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "The post has been deleted successfully",
	})
}
