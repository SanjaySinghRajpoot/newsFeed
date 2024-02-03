package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
)

func SendNotification(post models.Post) error {

	// get the user followers
	followingList := make([]models.Follower, 0)
	result := config.DB.Where("following_user_id = ?", post.UserID).Find(&followingList)

	if result.Error != nil {
		return result.Error
	}

	useIDs := make([]uint, 0)

	for _, user := range followingList {
		useIDs = append(useIDs, user.FollowerUserID)
	}

	// call the notification service
	url := "http://localhost:8081/user/notification"

	payload := models.NotificationPayload{
		UserID:      useIDs,
		Type:        "inapp",
		Description: post.Content,
	}

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		fmt.Println("Error while Marshalling:", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
