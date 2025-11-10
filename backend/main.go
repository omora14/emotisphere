package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"emotisphere/services"
	"emotisphere/utils"
	ws "emotisphere/websocket"
)

func main() {
	// environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// logger
	utils.InitLogger()

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()

	// processor
	processor := services.NewProcessor(hub)

	// routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get countries from query parameter or use default (max 5 for free tier)
		// Default to 5 countries: USA, Costa Rica, Brazil, Bolivia, Spain
		countries := []string{"us", "cr", "br", "bo", "es"}
		if countryParam := r.URL.Query().Get("countries"); countryParam != "" {
			// Parse comma-separated countries
			countries = []string{}
			parts := strings.Split(countryParam, ",")
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" && len(countries) < 5 {
					countries = append(countries, trimmed)
				}
			}
		}

		interval := 5 * time.Minute
		if intervalParam := r.URL.Query().Get("interval"); intervalParam != "" {
			if parsed, err := time.ParseDuration(intervalParam); err == nil {
				interval = parsed
			}
		}

		go processor.Start(interval, countries)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processor started"))
		utils.LogInfo("Processor started via API")
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		processor.Stop()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processor stopped"))
		utils.LogInfo("Processor stopped via API")
	})

	// Start processor automatically if API keys are set
	// Using 5 countries: USA, Costa Rica, Brazil, Bolivia, Spain
	if os.Getenv("NEWSDATA_API_KEY") != "" && os.Getenv("HUGGINGFACE_API_KEY") != "" {
		countries := []string{"us", "cr", "br", "bo", "es"} // 5 countries for free tier
		interval := 10 * time.Minute                        // Increased interval to avoid rate limits
		go processor.Start(interval, countries)
		utils.LogInfo("Processor started automatically with 5 countries: USA, Costa Rica, Brazil, Bolivia, Spain")
	} else {
		utils.LogInfo("API keys not set, processor will not start automatically. Use /start endpoint to start manually.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.LogInfo("Server starting on port %s", port)
	utils.LogInfo("WebSocket endpoint: ws://localhost:%s/ws", port)
	utils.LogInfo("Health check: http://localhost:%s/health", port)
	utils.LogInfo("Start processor: POST http://localhost:%s/start", port)
	utils.LogInfo("Stop processor: POST http://localhost:%s/stop", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		utils.LogError("Server failed to start: %v", err)
		log.Fatal(err)
	}
}

func handleWebSocket(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // all origins for development
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.LogError("WebSocket upgrade failed: %v", err)
		return
	}

	client := ws.NewClient(hub, conn)
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
