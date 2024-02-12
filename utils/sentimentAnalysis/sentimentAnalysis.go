package sentimentanalysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"gorm.io/datatypes"
)

func GetSentimentAnalysis(post models.Post) (models.SentimentAnalysisResponse, error) {
	url := "http://127.0.0.1:5000/analysis"

	payload := models.SentimentAnalysisPayload{
		Type:  "text",
		Value: post.Content,
	}

	jsonStr, err := json.Marshal(&payload)
	if err != nil {
		fmt.Println("Error while Marshalling:", err)
		return models.SentimentAnalysisResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return models.SentimentAnalysisResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.SentimentAnalysisResponse{}, err
	}
	defer resp.Body.Close()

	var snsAnalysis models.SentimentAnalysisResponse
	err = json.NewDecoder(resp.Body).Decode(&snsAnalysis)
	if err != nil {
		return models.SentimentAnalysisResponse{}, err
	}

	return snsAnalysis, nil
}

func UpdateSentimentAnalysisCRON() {

	var userPosts []models.Post
	result := config.DB.Debug().Where("created_at BETWEEN NOW() - INTERVAL '24 HOURS' AND NOW()").
		Where("sentiment_analysis->>'prominent_sentiment' = ?", "").
		Find(&userPosts)
	if result.Error != nil {
		log := fmt.Sprintf("Error unable to get the data: %s", result.Error)
		fmt.Println(log)
		return
	}

	for i, post := range userPosts {

		fmt.Println(i)
		fmt.Println(post.Content)

		snsAnalysis, err := GetSentimentAnalysis(post)
		if err != nil {
			log := fmt.Sprintf("Error unable to fetch the data from server: %s", err)
			fmt.Println(log)
			return
		}

		snsUpdate := models.SentimentAnalysis{
			ProminentSentiment: snsAnalysis.ProminentSentiment,
			ScorePositive:      snsAnalysis.ScorePositive,
			ScoreNegative:      snsAnalysis.ScoreNegative,
			ScoreNeutral:       snsAnalysis.ScoreNeutral,
		}

		updatePost := models.Post{
			SentimentAnalysis: datatypes.NewJSONType(snsUpdate),
		}

		// Update the post
		result = config.DB.Model(&post).Updates(&updatePost)

		if result.Error != nil {
			log := fmt.Sprintf("Error unable to save the data: %s", result.Error)
			fmt.Println(log)
			return
		}

	}

	return
}
