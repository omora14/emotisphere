# Emotisphere

**Real-Time World Emotion Map** - An interactive visualization platform that displays emotions from social posts and text messages on a live, global map.

## Project Overview

Emotisphere is a full-stack application that combines:
- **AI-Powered Emotion Analysis** using Hugging Face models
- **Real-Time Data Streaming** via WebSocket connections
- **Interactive Map Visualization** using React Leaflet
- **Go Backend** for high-performance emotion processing

The map displays emotions as colored markers that update in real-time, creating a "heatmap" of global emotional states similar to Google Photos but for human feelings.

## Project Structure

```
emotisphere/
├── frontend/              # React + Vite frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   │   ├── MapView.jsx        # Main map visualization
│   │   │   ├── emotionLegend.jsx  # Emotion legend component
│   │   │   └── realtimeFeed.jsx   # Real-time feed component
│   │   ├── services/      # API and WebSocket services
│   │   │   └── websocket.js       # WebSocket client
│   │   ├── mock/          # Mock data for development
│   │   │   └── data.json          # Sample emotion data
│   │   ├── App.jsx        # Main app component
│   │   └── main.jsx       # Entry point
│   ├── package.json       # Frontend dependencies
│   └── vite.config.js     # Vite configuration
│
├── backend/               # Go backend server
│   ├── main.go            # Server entry point
│   ├── websocket/         # WebSocket handlers
│   │   ├── hub.go         # Connection hub
│   │   ├── client.go      # Client management
│   │   └── message.go     # Message types
│   ├── analysis/          # Emotion analysis
│   │   └── emotion.go     # Hugging Face integration
│   ├── utils/             # Utilities
│   │   └── logger.go      # Logging utilities
│   ├── mock/              # Backend mock data
│   │   └── data.json      # Sample data
│   └── go.mod             # Go dependencies
│
└── README.md              # This file
```

## Getting Started

### Prerequisites

- **Node.js** (v18 or higher)
- **Go** (v1.21 or higher)
- **npm** or **yarn**

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173` (or the port shown in terminal).

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create a `.env` file (if needed for API keys):
```bash
# Add your Hugging Face API token if required
HUGGINGFACE_API_KEY=your_api_key_here
```

4. Run the backend server:
```bash
go run main.go
```

The backend will typically run on `http://localhost:8080`.

## Features

### Current Implementation
- Interactive world map with OpenStreetMap tiles
- Emotion markers with color coding:
  - 🟡 **Happy** - Gold/Yellow
  - 🔵 **Sad** - Royal Blue
  - 🔴 **Angry** - Crimson Red
  - 🟠 **Surprised** - Dark Orange
- Mock data visualization
- Popup details showing emotion, intensity, and coordinates

### Planned Features
- [ ] Real-time WebSocket connection
- [ ] Hugging Face emotion analysis integration
- [ ] Live social media post processing
- [ ] Emotion intensity heatmap overlay
- [ ] Time-based emotion filtering
- [ ] Emotion statistics dashboard

## Technology Stack

### Frontend
- **React 19** - UI framework
- **Vite** - Build tool and dev server
- **React Leaflet** - Map visualization
- **Tailwind CSS** - Styling (configured)

### Backend
- **Go 1.25** - Server language
- **Gorilla WebSocket** - WebSocket support
- **Hugging Face API** - Emotion analysis (planned)

## Development Notes

### Starting with Mock Data
The project currently uses mock data located in `frontend/src/mock/data.json`. This allows for development and testing without requiring live data sources or API integrations.

### Next Steps
1. Set up WebSocket connection between frontend and backend
2. Integrate Hugging Face emotion analysis API
3. Create data ingestion pipeline for social posts
4. Implement real-time updates on the map
5. Add filtering and statistics features

## License

This project is part of a learning sprint and is intended for portfolio/resume purposes.

## Contributing

This is a personal learning project, but suggestions and feedback are welcome!

---

**Sprint Goal**: Successfully build and connect all parts of the Real-Time World Emotion Map, combining AI, real-time data, and interactive visuals.

