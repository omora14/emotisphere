package services

import (
	"log"
	"time"

	"emotisphere/websocket"
)

// Processor handles the emotion analysis pipeline
type Processor struct {
	NewsService     *NewsService
	EmotionService  *EmotionService
	LocationService *LocationService
	Hub             *websocket.Hub
	Running         bool
	StopChan        chan bool
}

// NewProcessor creates a new processor
func NewProcessor(hub *websocket.Hub) *Processor {
	return &Processor{
		NewsService:     NewNewsService(),
		EmotionService:  NewEmotionService(),
		LocationService: NewLocationService(),
		Hub:             hub,
		StopChan:        make(chan bool),
	}
}

func (p *Processor) Start(interval time.Duration, countries []string) {
	if p.Running {
		log.Println("Processor is already running")
		return
	}

	p.Running = true
	log.Printf("Starting processor with interval: %v, countries: %v", interval, countries)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go p.ProcessBatch(countries)

	for {
		select {
		case <-ticker.C:
			go p.ProcessBatch(countries)
		case <-p.StopChan:
			log.Println("Stopping processor")
			p.Running = false
			return
		}
	}
}

func (p *Processor) Stop() {
	if p.Running {
		p.StopChan <- true
	}
}

func (p *Processor) ProcessBatch(countries []string) {
	log.Printf("Fetching news articles for countries: %v", countries)

	articles, err := p.NewsService.FetchNews(countries)
	if err != nil {
		log.Printf("Error fetching news: %v", err)
		// Don't broadcast error, just log it - try processing individual countries
		// Try fetching for each country individually as fallback
		for _, country := range countries {
			countryArticles, countryErr := p.NewsService.FetchNews([]string{country})
			if countryErr == nil && len(countryArticles) > 0 {
				log.Printf("Successfully fetched %d articles for %s", len(countryArticles), country)
				articles = append(articles, countryArticles...)
			} else {
				log.Printf("Could not fetch articles for %s: %v", country, countryErr)
			}
		}

		if len(articles) == 0 {
			log.Printf("No articles fetched for any country")
			return
		}
	}

	log.Printf("Fetched %d articles total, processing...", len(articles))

	for i, article := range articles {
		go func(art NewsArticle, index int) {
			if index > 0 {
				time.Sleep(time.Second * 2) // 2 second delay between requests
			}
			p.ProcessArticle(art)
		}(article, i)
	}
}

func (p *Processor) ProcessArticle(article NewsArticle) {
	text := p.NewsService.ExtractText(article)
	if text == "" {
		return
	}

	// Analyze
	emotion, intensity, err := p.EmotionService.AnalyzeEmotion(text)
	if err != nil {
		log.Printf("Error analyzing emotion: %v", err)
		return
	}

	city, country, err := p.LocationService.ProcessLocation(article.Country)
	if err != nil {
		log.Printf("Error processing location: %v", err)
		return
	}

	// coordinates
	lat, lng, err := p.LocationService.GetCoordinates(city, country)
	if err != nil {
		log.Printf("Error getting coordinates for %s, %s: %v", city, country, err)
		return
	}

	// emotion data
	emotionData := websocket.EmotionData{
		City:      city,
		Country:   country,
		Emotion:   emotion,
		Intensity: intensity,
		Lat:       lat,
		Lng:       lng,
		Text:      text[:min(100, len(text))],
	}

	p.Hub.Broadcast <- websocket.Message{
		Type: websocket.MessageTypeEmotion,
		Data: emotionData,
	}

	log.Printf("Processed: %s - %.2f at %s, %s", emotion, intensity, city, country)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
