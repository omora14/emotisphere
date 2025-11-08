import { useEffect } from 'react';
import { MapContainer, TileLayer, CircleMarker, Popup } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import mockData from '../mock/data.json';

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
  useEffect(() => {
    console.log('MapView loaded with data:', mockData);
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

        {mockData.map((item, idx) => (
          <CircleMarker
            key={idx}
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
                <span style={{ fontSize: '12px', color: '#666' }}>
                  {item.lat.toFixed(4)}, {item.lng.toFixed(4)}
                </span>
              </div>
            </Popup>
          </CircleMarker>
        ))}
      </MapContainer>
    </div>
  );
}
