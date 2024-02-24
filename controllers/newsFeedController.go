package controllers

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
	"github.com/gin-gonic/gin"
)

//  we need to create a precomputed cache for a given user and store it in cache

func GetNewsFeed(c *gin.Context) {

	userID := helpers.GetAuthUser(c).ID

	// get all the people who this is following
	followList := make([]models.Follower, 0)
	result := config.DB.Where("follower_user_id = ?", userID).Find(&followList)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	posts, err := redis.GetPostCache(userID)
	if err != nil {
		log.Fatal("Post Cache Not found in getNewsFeed")
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

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(posts), func(i, j int) { posts[i], posts[j] = posts[j], posts[i] })

	c.JSON(http.StatusOK, gin.H{
		"post": totalPosts,
	})

}
