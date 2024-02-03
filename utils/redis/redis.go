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
