package controllers

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/validations.go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CommentOnPost comments on a post
func CommentOnPost(c *gin.Context) {
	// Get data from request
	var userInput struct {
		PostId uint   `json:"postId" binding:"required,min=1"`
		Body   string `json:"body" binding:"required,min=1"`
	}

	err := c.ShouldBindJSON(&userInput)

	// Validate the data
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errors),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !validations.IsExistValue("posts", "id", userInput.PostId) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"PostId": "The post does not exist",
			},
		})
		return
	}

	// Store the comment
	authId := helpers.GetAuthUser(c).ID

	comment := models.Comment{
		PostID: userInput.PostId,
		UserID: authId,
		Body:   userInput.Body,
	}

	result := config.DB.Create(&comment)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// Return the comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// EditComment finds a comment by id
func GetComment(c *gin.Context) {
	// Get the comment id from url
	id := c.Param("comment_id")

	// Find the comment
	var comment models.Comment
	result := config.DB.First(&comment, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Return the comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// UpdateComment updates a comment of a post
func UpdateComment(c *gin.Context) {
	// Get the id from url
	id := c.Param("comment_id")

	var userInput struct {
		Body string `json:"body" binding:"required,min=1"`
	}

	// Validate in request
	err := c.ShouldBindJSON(&userInput)

	// Validate the data
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errors),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find the comment
	var comment models.Comment
	result := config.DB.First(&comment, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Update the comment
	comment.Body = userInput.Body
	result = config.DB.Save(&comment)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// Return the comment
	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

func DeleteComment(c *gin.Context) {
	// Get the id from url
	id := c.Param("comment_id")

	// Find the comment
	var comment models.Comment
	result := config.DB.First(&comment, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Delete the comment
	config.DB.Delete(&comment)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "The comment has been deleted successfully!",
	})
}
