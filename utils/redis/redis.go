package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
		// use os.env
		Addr:     "localhost:6379",
		Password: password,
		DB:       0,
	})

}

func SetRedisData(email string, tokenString string) (string, error) {

	ctx := context.Background()

	err := RedisClient.Set(ctx, email, tokenString, 24*time.Hour).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetRedisData(email string) (string, error) {

	ctx := context.Background()

	tokenString, err := RedisClient.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetUserID(UserID string, Count int) (string, error) {

	ctx := context.Background()

	err := RedisClient.Set(ctx, UserID, Count, 30*time.Minute).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetUserID(UserID string) (int, error) {

	ctx := context.Background()

	cnt, err := RedisClient.Get(ctx, UserID).Int()

	if err != nil {
		return -1, err
	}

	return cnt, nil
}

func SetPostCache(UserID uint, post models.Post) (string, error) {
	ctx := context.Background()

	redisKey := fmt.Sprintf("%d", UserID)

	currentTime := time.Now()

	postTimeStr := currentTime.Format(time.RFC3339Nano)

	postBytes, err := json.Marshal(&post)
	if err != nil {
		return "", err
	}

	if RedisClient == nil {
		fmt.Println("redis not working")
		return "", nil
	}

	err = RedisClient.HSet(ctx, redisKey, postTimeStr, string(postBytes)).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetPostCache(UserID uint) ([]models.Post, error) {

	ctx := context.Background()

	userIDstr := fmt.Sprintf("%d", UserID)

	var allPosts []models.Post

	str, err := RedisClient.HGetAll(ctx, userIDstr).Result()

	if err != nil {
		return []models.Post{}, err
	}

	fmt.Println(str)

	for _, val := range str {
		var post models.Post

		err = json.Unmarshal([]byte(val), &post)
		if err != nil {
			fmt.Println("Error:", err)
			return allPosts, nil
		}

		allPosts = append(allPosts, post)

	}

	return allPosts, nil
}
