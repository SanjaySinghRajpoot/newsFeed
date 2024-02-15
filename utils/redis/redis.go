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
		Addr:     "cache:6379",
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

	postBytes, err := json.Marshal(&post)
	if err != nil {
		return "", err
	}

	err = RedisClient.Set(ctx, redisKey, postBytes, 30*time.Minute).Err()

	if err != nil {
		return "Something went wrong", err
	}

	return "", nil
}

func GetPostCache(UserID uint) (models.Post, error) {

	ctx := context.Background()

	userIDstr := fmt.Sprintf("%d", UserID)

	var post models.Post

	str, err := RedisClient.Get(ctx, userIDstr).Bytes()

	if err != nil {
		return models.Post{}, err
	}

	err = json.Unmarshal(str, &post)

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}
