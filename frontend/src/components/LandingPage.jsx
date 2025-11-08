import { useState } from 'react';
import './LandingPage.css';

export default function LandingPage({ onEnterMap }) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <div className="landing-page">
      <div className="landing-content">
        <div className="landing-header">
          <h1 className="landing-title">Emotisphere</h1>
          <div className="landing-subtitle">Real-Time World Emotion Map</div>
        </div>

        <div className="landing-description">
          <p className="description-text">
            Explore the emotional pulse of our planet in real-time. 
            Emotisphere visualizes emotions from around the world on an interactive map, 
            creating a living tapestry of human feelings.
          </p>
          
          <div className="features-grid">
            <div className="feature-item">
              <div className="feature-icon">🌍</div>
              <div className="feature-text">Global Visualization</div>
            </div>
            <div className="feature-item">
              <div className="feature-icon">⚡</div>
              <div className="feature-text">Real-Time Updates</div>
            </div>
            <div className="feature-item">
              <div className="feature-icon">🤖</div>
              <div className="feature-text">AI-Powered Analysis</div>
            </div>
          </div>
        </div>

        <button 
          className="enter-button"
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
          onClick={onEnterMap}
        >
          <span className="button-text">Explore the Map</span>
          <span className={`button-arrow ${isHovered ? 'hovered' : ''}`}>→</span>
        </button>
      </div>
    </div>
  );
}

