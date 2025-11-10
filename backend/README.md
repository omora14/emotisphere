# Emotisphere Backend

Real-time emotion analysis backend that fetches news articles, analyzes emotions, and broadcasts results via WebSocket.

## Architecture

The backend follows this workflow:

1. **News Fetching** → Fetches recent news articles from newsdata.io
2. **Emotion Analysis** → Analyzes text using Hugging Face emotion model
3. **Location Mapping** → Maps country/city to coordinates using Nominatim
4. **WebSocket Broadcast** → Sends emotion data to connected clients in real-time

## Prerequisites

- Go 1.21 or higher
- API Keys:
  - [NewsData.io API Key](https://newsdata.io/) (free tier available)
  - [Hugging Face API Key](https://huggingface.co/settings/tokens) (free tier available)

## Setup

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` and add your API keys:

```env
NEWSDATA_API_KEY=your_newsdata_api_key_here
HUGGINGFACE_API_KEY=your_huggingface_api_key_here
PORT=8080
```

### 3. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env`).

## API Endpoints

### WebSocket Connection

**Endpoint:** `ws://localhost:8080/ws`

Connect to this endpoint to receive real-time emotion data. The server will send messages in this format:

```json
{
  "type": "emotion",
  "data": {
    "city": "",
    "country": "United States",
    "emotion": "happy",
    "intensity": 0.85,
    "lat": 39.8283,
    "lng": -98.5795,
    "text": "Sample text..."
  }
}
```

### Health Check

**GET** `http://localhost:8080/health`

Returns `200 OK` if the server is running.

### Start Processor

**POST** `http://localhost:8080/start?countries=us,gb,jp&interval=5m`

Starts the emotion analysis processor.

**Query Parameters:**
- `countries` (optional): Comma-separated country codes (default: us,gb,jp,ca,au,de,fr)
- `interval` (optional): Processing interval in Go duration format (default: 5m)

**Example:**
```bash
curl -X POST "http://localhost:8080/start?countries=us,gb&interval=3m"
```

### Stop Processor

**POST** `http://localhost:8080/stop`

Stops the emotion analysis processor.

**Example:**
```bash
curl -X POST http://localhost:8080/stop
```

## Testing

### Test with Mock Data (No API Keys Required)

If you don't have API keys yet, you can test the WebSocket connection:

1. Start the server:
```bash
go run main.go
```

2. In another terminal, use a WebSocket client to connect:
```bash
# Using wscat (install with: npm install -g wscat)
wscat -c ws://localhost:8080/ws
```

3. The server will connect, but won't send data until API keys are configured.

### Test with API Keys

1. Add your API keys to `.env`
2. Start the server - it will automatically start processing
3. Connect a WebSocket client to see real-time emotion data

### Manual Testing

You can test individual services:

```go
// Test news fetching
newsService := services.NewNewsService()
articles, err := newsService.FetchNews([]string{"us", "gb"})

// Test emotion analysis
emotionService := services.NewEmotionService()
emotion, intensity, err := emotionService.AnalyzeEmotion("I'm so happy today!")

// Test location mapping
locationService := services.NewLocationService()
lat, lng, err := locationService.GetCoordinates("", "United States")
```

## Project Structure

```
backend/
├── main.go              # Server entry point
├── services/
│   ├── news.go          # NewsData.io integration
│   ├── emotion.go       # Hugging Face emotion analysis
│   ├── location.go      # Location to coordinates mapping
│   └── processor.go     # Main processing pipeline
├── websocket/
│   ├── hub.go           # WebSocket hub
│   ├── client.go        # WebSocket client
│   └── message.go      # Message types
├── utils/
│   └── logger.go        # Logging utilities
├── .env.example         # Environment variables template
└── README.md            # This file
```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `NEWSDATA_API_KEY` | NewsData.io API key | Yes | - |
| `HUGGINGFACE_API_KEY` | Hugging Face API key | Yes | - |
| `HUGGINGFACE_MODEL` | Hugging Face model name | No | `j-hartmann/emotion-english-distilroberta-base` |
| `PORT` | Server port | No | `8080` |

### Supported Countries

The processor supports ISO country codes. Common ones:
- `us` - United States
- `gb` - United Kingdom
- `jp` - Japan
- `ca` - Canada
- `au` - Australia
- `de` - Germany
- `fr` - France
- And many more...

## Troubleshooting

### "API key not set" errors

Make sure your `.env` file exists and contains valid API keys. The server will not start processing automatically if keys are missing.

### WebSocket connection fails

- Check that the server is running on the correct port
- Verify CORS settings if connecting from a browser
- Check browser console for connection errors

### No emotion data received

- Verify API keys are correct
- Check server logs for errors
- Ensure the processor is running (check `/health` endpoint)
- Verify news articles are being fetched successfully

### Rate limiting

Both NewsData.io and Hugging Face have rate limits on free tiers:
- **NewsData.io**: 200 requests/day (free tier)
- **Hugging Face**: 1000 requests/month (free tier)

Adjust the processing interval if you hit rate limits.

## Notes

- The processor runs in the background and processes articles at regular intervals
- Each article is processed asynchronously
- Location mapping uses Nominatim (OpenStreetMap) which is free and doesn't require an API key
- The emotion model can be changed via `HUGGINGFACE_MODEL` environment variable
- Processing interval can be adjusted via the `/start` endpoint

## API Documentation

- [NewsData.io API](https://newsdata.io/documentation)
- [Hugging Face Inference API](https://huggingface.co/docs/api-inference/index)
- [Nominatim API](https://nominatim.org/release-docs/develop/api/Overview/)

