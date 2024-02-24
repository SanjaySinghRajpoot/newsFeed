package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/redis"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	NEWSFEED = "newsfeed"
)

func main() {

	// Connect to the postgres database
	config.Connect()

	// Redis Cache Setup
	redis.RedisClient = redis.SetUpRedis("12345678")

	// Set up configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Replace with your Kafka broker address
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	// Create consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Subscribe to a topics
	topics := []string{"newsfeed"}
	consumer.SubscribeTopics(topics, nil)

	// Handle messages and shutdown signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	var postObj models.Post
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		default:
			// time out of 100 millisecond
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:

				err = json.Unmarshal(e.Value, &postObj)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, postObj.Content)
				messageType := e.TopicPartition.Topic
				switch *messageType {
				case NEWSFEED:
					handleNewsfeed(e, postObj)
				}

			case kafka.Error:
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				run = false

			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func handleNewsfeed(e *kafka.Message, post models.Post) {

	fmt.Println("working newsfeed push model")

	followList := make([]models.Follower, 0)
	result := config.DB.Where("following_user_id = ?", post.UserID).Find(&followList)

	if result.Error != nil {
		fmt.Println("Error while Getting the follower list:", result.Error)
		return
	}

	for _, friend := range followList {

		fmt.Println(friend.FollowingUserID)
		fmt.Println(friend.FollowerUserID)

		msg, er := redis.SetPostCache(friend.FollowerUserID, post)
		if er != nil {
			fmt.Println(msg)
			// formatError.InternalServerError(c, er)
			return
		}

	}

	// url := "http://localhost:8082/sms"

	// payload := models.ServicePayload{
	// 	Notification_id: notifID,
	// 	UserID:          userID,
	// 	Message:         string(e.Value),
	// }

	// jsonStr, err := json.Marshal(&payload)
	// if err != nil {
	// 	fmt.Println("Error while Marshalling:", err)
	// 	return
	// }

	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

}
