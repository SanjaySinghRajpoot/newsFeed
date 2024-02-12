package controllers

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
	"github.com/gin-gonic/gin"
)

func GetNewsFeed(c *gin.Context) {

	userID := helpers.GetAuthUser(c).ID

	// get all the people who this is following

	followList := make([]models.Follower, 0)
	result := config.DB.Where("follower_user_id = ?", userID).Find(&followList)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// check if the follower has posted anything in the past 24 hours
	posts := make([]models.Post, 0)

	for _, follow := range followList {

		var userPosts models.Post
		userPosts, err := redis.GetPostCache(follow.FollowingUserID)
		if err != nil {
			formatError.InternalServerError(c, err)
			return
		}

		posts = append(posts, userPosts)

		// var getDBposts models.Post
		// result := config.DB.Where("user_id = ? AND created_at BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()", follow.FollowingUserID).Find(&userPosts)

		// if result.Error != nil {
		// 	formatError.InternalServerError(c, result.Error)
		// 	return
		// }
	}

	c.JSON(http.StatusOK, gin.H{
		"post": posts,
	})

}
