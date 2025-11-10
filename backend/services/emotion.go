package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type EmotionResponse struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

type HuggingFaceResponse []EmotionResponse

type EmotionService struct {
	APIKey string
	Client *http.Client
	Model  string
}

func NewEmotionService() *EmotionService {
	model := os.Getenv("HUGGINGFACE_MODEL")
	if model == "" {
		model = "j-hartmann/emotion-english-distilroberta-base"
	}

	return &EmotionService{
		APIKey: os.Getenv("HUGGINGFACE_API_KEY"),
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Model: model,
	}
}

func (es *EmotionService) AnalyzeEmotion(text string) (string, float64, error) {
	if es.APIKey == "" {
		return "", 0, fmt.Errorf("HUGGINGFACE_API_KEY not set in environment variables")
	}

	endpoints := []string{
		fmt.Sprintf("https://router.huggingface.co/hf-inference/v1/models/%s", es.Model),
		fmt.Sprintf("https://router.huggingface.co/hf-inference/models/%s", es.Model),
		fmt.Sprintf("https://api-inference.huggingface.co/models/%s", es.Model), // Fallback to old endpoint
	}

	var lastErr error
	for _, url := range endpoints {
		result, err := es.tryAnalyzeWithEndpoint(url, text)
		if err == nil {
			return result.emotion, result.score, nil
		}
		lastErr = err
		if err != nil && (strings.Contains(err.Error(), "410") || strings.Contains(err.Error(), "404")) {
			continue
		}
	}

	return "", 0, fmt.Errorf("all endpoint attempts failed, last error: %w", lastErr)
}

type emotionResult struct {
	emotion string
	score   float64
}

// tryAnalyzeWithEndpoint tries to analyze emotion with a specific endpoint according to the documentation at https://huggingface.co/docs/api-inference/quicktour
func (es *EmotionService) tryAnalyzeWithEndpoint(url, text string) (*emotionResult, error) {
	payload := map[string]interface{}{
		"inputs": text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", es.APIKey))
	req.Header.Set("x-use-cache", "false")

	// Make request
	resp, err := es.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// it's a 503 (model loading)
		if resp.StatusCode == http.StatusServiceUnavailable {
			time.Sleep(5 * time.Second)
			return nil, fmt.Errorf("model is loading (503)")
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response - Handle different response formats
	emotions, err := es.parseEmotionResponse(body)
	if err != nil {
		return nil, err
	}

	if len(emotions) == 0 {
		return nil, fmt.Errorf("no emotions found in response")
	}

	bestEmotion := emotions[0]
	for _, emotion := range emotions {
		if emotion.Score > bestEmotion.Score {
			bestEmotion = emotion
		}
	}

	emotionLabel := mapEmotionLabel(bestEmotion.Label)

	return &emotionResult{emotion: emotionLabel, score: bestEmotion.Score}, nil
}

// parseEmotionResponse parses the emotion response from Hugging Face API
func (es *EmotionService) parseEmotionResponse(body []byte) (HuggingFaceResponse, error) {
	var emotions HuggingFaceResponse

	if err := json.Unmarshal(body, &emotions); err == nil && len(emotions) > 0 {
		return emotions, nil
	}

	var responseData interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if emotionArray, ok := responseData.([]interface{}); ok {
		for _, item := range emotionArray {
			if emotionMap, ok := item.(map[string]interface{}); ok {
				if label, ok := emotionMap["label"].(string); ok {
					var score float64
					if scoreVal, ok := emotionMap["score"].(float64); ok {
						score = scoreVal
					}
					emotions = append(emotions, EmotionResponse{Label: label, Score: score})
				}
			}
		}
		if len(emotions) > 0 {
			return emotions, nil
		}
	}

	if emotionMap, ok := responseData.(map[string]interface{}); ok {
		if data, ok := emotionMap["data"].([]interface{}); ok {
			for _, item := range data {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if label, ok := itemMap["label"].(string); ok {
						var score float64
						if scoreVal, ok := itemMap["score"].(float64); ok {
							score = scoreVal
						}
						emotions = append(emotions, EmotionResponse{Label: label, Score: score})
					}
				}
			}
		} else if label, ok := emotionMap["label"].(string); ok {
			var score float64
			if scoreVal, ok := emotionMap["score"].(float64); ok {
				score = scoreVal
			}
			emotions = append(emotions, EmotionResponse{Label: label, Score: score})
		}
	}

	if len(emotions) == 0 {
		return nil, fmt.Errorf("could not parse emotion response, unexpected format")
	}

	return emotions, nil
}

func mapEmotionLabel(hfLabel string) string {
	labelMap := map[string]string{
		"joy":      "happy",
		"sadness":  "sad",
		"anger":    "angry",
		"fear":     "surprised", // fear to surprised for now
		"surprise": "surprised",
		"love":     "happy",
		"neutral":  "neutral",
	}

	if mapped, ok := labelMap[hfLabel]; ok {
		return mapped
	}
	return "neutral"
}
