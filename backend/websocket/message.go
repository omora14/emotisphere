package websocket

// represents a WebSocket message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// EmotionData represents emotion data with location
type EmotionData struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Emotion   string  `json:"emotion"`
	Intensity float64 `json:"intensity"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Text      string  `json:"text,omitempty"` // Optional: original text
}

// Message types
const (
	MessageTypeEmotion = "emotion"
	MessageTypeError   = "error"
	MessageTypeInfo    = "info"
)
