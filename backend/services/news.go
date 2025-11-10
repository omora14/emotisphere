package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// represents a news article from newsdata.io
type NewsArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Country     []string `json:"country"`
	Language    string   `json:"language"`
	PubDate     string   `json:"pubDate"`
}

// represents the response from newsdata.io API
type NewsResponse struct {
	Status   string        `json:"status"`
	Results  []NewsArticle `json:"results"`
	NextPage string        `json:"nextPage,omitempty"`
}

// NewsService handles fetching news data
type NewsService struct {
	APIKey string
	Client *http.Client
}

// NewNewsService creates a new news service
func NewNewsService() *NewsService {
	return &NewsService{
		APIKey: os.Getenv("NEWSDATA_API_KEY"),
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchNews fetches recent news articles
func (ns *NewsService) FetchNews(countries []string) ([]NewsArticle, error) {
	if ns.APIKey == "" {
		return nil, fmt.Errorf("NEWSDATA_API_KEY not set in environment variables")
	}

	if len(countries) > 5 {
		countries = countries[:5]
	}

	url := fmt.Sprintf("https://newsdata.io/api/1/news?apikey=%s&language=en", ns.APIKey)

	if len(countries) > 0 {
		countryStr := ""
		for i, country := range countries {
			if i > 0 {
				countryStr += ","
			}
			countryStr += country
		}
		url += fmt.Sprintf("&country=%s", countryStr)
	}

	// Add category for better results (optional, helps with free tier)
	// Using "top" category which is available in free tier
	url += "&category=top"

	resp, err := ns.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news API returned status %d: %s", resp.StatusCode, string(body))
	}

	var newsResponse NewsResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse news response: %w", err)
	}

	if newsResponse.Status != "success" {
		// If status is not success, return empty results instead of error
		// This allows processing to continue with other countries
		return []NewsArticle{}, nil
	}

	return newsResponse.Results, nil
}

func (ns *NewsService) ExtractText(article NewsArticle) string {
	// Prefer content, then description, then title
	if article.Content != "" {
		return article.Content
	}
	if article.Description != "" {
		return article.Description
	}
	return article.Title
}
