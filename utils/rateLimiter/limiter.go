package limiter

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
)

const (
	maxRequests     = 2
	perMinutePeriod = 1 * time.Minute
)

var (
	mutex = &sync.Mutex{}
)

func RateLimiter(userID uint) error {
	mutex.Lock()
	defer mutex.Unlock()

	userIDstr := fmt.Sprintf("%d", userID)

	// get from redis cache
	count, err := redis.GetUserID(userIDstr)

	if err != nil {
		fmt.Printf("Failed to Get the Redis Cache")
		fmt.Println(err)
	}

	if count >= maxRequests {
		errReq := errors.New("too many requests")
		return errReq
	}

	count = count + 1

	msg, err := redis.SetUserID(userIDstr, count)
	if err != nil {
		fmt.Printf("Failed to set the Redis Cache: %s", msg)

		fmt.Println(err)

		return err
	}

	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		count, err := redis.GetUserID(userIDstr)

		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache: %d", count)

			return
		}

		count = count - 1

		msg, err := redis.SetUserID(userIDstr, count)
		if err != nil {
			fmt.Printf("Failed to Set the Redis Cache: %s", msg)

			return
		}

	})

	return nil
}
