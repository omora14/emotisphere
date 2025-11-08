import { useState } from 'react';
import './App.css';
import 'leaflet/dist/leaflet.css';
import MapView from './components/MapView';
import LandingPage from './components/LandingPage';
import Header from './components/Header';
import Footer from './components/Footer';

function App() {
  const [showMap, setShowMap] = useState(false);

  const handleEnterMap = () => {
    setShowMap(true);
  };

  const handleBackToLanding = () => {
    setShowMap(false);
  };

  return (
    <div className="app-container">
      {showMap ? (
        <>
          <Header onBackToLanding={handleBackToLanding} />
          <MapView />
          <Footer />
        </>
      ) : (
        <LandingPage onEnterMap={handleEnterMap} />
      )}
    </div>
  );
}

export default App;
