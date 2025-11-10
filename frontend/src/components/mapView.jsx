import { useEffect, useState, useRef } from 'react';
import { MapContainer, TileLayer, CircleMarker, Popup } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import websocketService from '../services/websocket';

// for default marker icons in react-leaflet this is the best way to fix it apparently
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
});

// Emotion color mapping, please check below to see the colors and map them
const getEmotionColor = (emotion) => {
  const colorMap = {
    happy: '#FFD700',      // Gold/Yellow
    sad: '#4169E1',        // Royal Blue
    angry: '#DC143C',      // Crimson Red
    surprised: '#FF8C00',  // Dark Orange
    neutral: '#808080',    // Gray
  };
  return colorMap[emotion] || colorMap.neutral;
};

export default function MapView() {
  const [emotionData, setEmotionData] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const dataMapRef = useRef(new Map()); // Use Map to store unique entries by location+emotion

  useEffect(() => {
    // Connect to WebSocket
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
    websocketService.connect(wsUrl);

    // Set up event listeners
    const handleConnect = () => {
      console.log('Connected to WebSocket');
      setIsConnected(true);
    };

    const handleDisconnect = () => {
      console.log('Disconnected from WebSocket');
      setIsConnected(false);
    };

    const handleEmotion = (data) => {
      console.log('Received emotion data:', data);
      
      // unique key for this location+emotion combination
      const key = `${data.lat}_${data.lng}_${data.emotion}`;
      
      dataMapRef.current.set(key, {
        ...data,
        id: key, 
        timestamp: Date.now(), 
      });

      setEmotionData(Array.from(dataMapRef.current.values()));
    };

    const handleError = (error) => {
      console.error('WebSocket error:', error);
    };

    websocketService.on('connect', handleConnect);
    websocketService.on('disconnect', handleDisconnect);
    websocketService.on('emotion', handleEmotion);
    websocketService.on('error', handleError);

    return () => {
      websocketService.off('connect', handleConnect);
      websocketService.off('disconnect', handleDisconnect);
      websocketService.off('emotion', handleEmotion);
      websocketService.off('error', handleError);
      websocketService.disconnect();
    };
  }, []);

  // here I am just setting the map bounds to prevent infinite scrolling
  const maxBounds = [
    [-85, -180], 
    [85, 180]   
  ];
  const minZoom = 2;
  const maxZoom = 18;

  return (
    <div style={{ 
      height: 'calc(100vh - 130px)', 
      width: '100%', 
      position: 'fixed',
      top: '70px',
      bottom: '60px',
      left: 0,
      right: 0,
      padding: 0,
      zIndex: 0
    }}>
      <MapContainer 
        center={[20, 0]} 
        zoom={2} 
        style={{ height: '100%', width: '100%' }}
        scrollWheelZoom={true}
        maxBounds={maxBounds}
        minZoom={minZoom}
        maxZoom={maxZoom}
        maxBoundsViscosity={1.0}
      >
        <TileLayer 
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" 
        />

        {/* Connection status indicator */}
        {!isConnected && (
          <div style={{
            position: 'absolute',
            top: '10px',
            left: '50%',
            transform: 'translateX(-50%)',
            backgroundColor: 'rgba(255, 0, 0, 0.8)',
            color: 'white',
            padding: '8px 16px',
            borderRadius: '4px',
            zIndex: 1000,
            fontSize: '14px',
          }}>
            Connecting to server...
          </div>
        )}

        {/* Emotion markers */}
        {emotionData.map((item) => (
          <CircleMarker
            key={item.id || `${item.lat}_${item.lng}_${item.emotion}`}
            center={[item.lat, item.lng]}
            radius={8 + (item.intensity * 12)}
            fillOpacity={0.6}
            color={getEmotionColor(item.emotion)}
            weight={2}
          >
            <Popup>
              <div style={{ textAlign: 'center' }}>
                <strong style={{ textTransform: 'capitalize', fontSize: '16px' }}>
                  {item.emotion}
                </strong>
                <br />
                <span style={{ fontSize: '14px' }}>
                  Intensity: {(item.intensity * 100).toFixed(0)}%
                </span>
                <br />
                {item.city && (
                  <>
                    <span style={{ fontSize: '12px', fontWeight: 'bold' }}>
                      {item.city}
                      {item.country && `, ${item.country}`}
                    </span>
                    <br />
                  </>
                )}
                <span style={{ fontSize: '12px', color: '#666' }}>
                  {item.lat.toFixed(4)}, {item.lng.toFixed(4)}
                </span>
                {item.text && (
                  <>
                    <br />
                    <span style={{ fontSize: '11px', color: '#888', fontStyle: 'italic', marginTop: '4px', display: 'block' }}>
                      "{item.text.substring(0, 50)}..."
                    </span>
                  </>
                )}
              </div>
            </Popup>
          </CircleMarker>
        ))}
      </MapContainer>
    </div>
  );
}
