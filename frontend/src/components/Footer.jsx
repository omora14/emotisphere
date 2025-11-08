import './Footer.css';

export default function Footer() {
  return (
    <footer className="app-footer">
      <div className="footer-content">
        <div className="footer-text">
          <span>Emotisphere</span>
          <span className="footer-separator">•</span>
          <span>Real-Time World Emotion Map</span>
        </div>
        <div className="footer-links">
          <a 
            href="https://github.com/omora14/emotisphere" 
            target="_blank" 
            rel="noopener noreferrer"
            className="footer-link"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  );
}

