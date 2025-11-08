import './Header.css';

export default function Header({ onBackToLanding }) {
  return (
    <header className="app-header">
      <div className="header-content">
        <div className="header-logo" onClick={onBackToLanding}>
          <span className="logo-text">Emotisphere</span>
        </div>
        <nav className="header-nav">
          <button className="nav-button" onClick={onBackToLanding}>
            About
          </button>
        </nav>
      </div>
    </header>
  );
}

