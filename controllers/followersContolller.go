package controllers

import (
	"net/http"
	"strconv"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/gin-gonic/gin"
)

func FollowRequest(c *gin.Context) {
	// Get the id from the url
	followingUserID := c.Param("id")

	followingInt, _ := strconv.Atoi(followingUserID)

	authID := helpers.GetAuthUser(c).ID

	follow := models.Follower{
		FollowingUserID: uint(followingInt),
		FollowerUserID:  authID,
	}

	result := config.DB.Create(&follow)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// Return the post
	c.JSON(http.StatusOK, gin.H{
		"message": follow,
	})
}

func UnfollowRequest(c *gin.Context) {
	// Get the id from the url
	followingUserID := c.Param("id")

	followingUserIDInt, _ := strconv.Atoi(followingUserID)

	followerUserID := helpers.GetAuthUser(c).ID

	user := models.Follower{
		FollowerUserID:  followerUserID,
		FollowingUserID: uint(followingUserIDInt),
	}

	// Delete the user
	result := config.DB.Delete(&user)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been unfollowed",
	})
}
