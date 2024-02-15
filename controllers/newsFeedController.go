package controllers

import (
	"log"
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

		// We are only getting the latest post published by the Following
		var userPosts models.Post
		userPosts, err := redis.GetPostCache(follow.FollowingUserID)
		if err != nil {
			log.Fatal("Post Cache Not found in getNewsFeed")
		}

		posts = append(posts, userPosts)

	}

	// Get the list of random posts from the DB in the past 24 hours with
	// Positive sentiment
	var getPositivePosts []models.Post
	result = config.DB.Debug().Where("created_at BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()").
		Where("sentiment_analysis->>'prominent_sentiment' = ?", "POSITIVE").Limit(10).Find(&getPositivePosts)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	var totalPosts []models.Post

	totalPosts = append(totalPosts, posts...)

	// Check for Unique Posts
	for _, friendPost := range posts {

		for _, positivePost := range getPositivePosts {

			if friendPost.ID != positivePost.ID {
				totalPosts = append(totalPosts, positivePost)
			}
		}
	}

	// rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })

	c.JSON(http.StatusOK, gin.H{
		"post": totalPosts,
	})

}
