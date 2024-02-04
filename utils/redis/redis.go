package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func SetUpRedis(password string) *redis.Client {

	return redis.NewClient(&redis.Options{
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
