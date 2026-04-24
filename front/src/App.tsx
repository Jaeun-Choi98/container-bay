import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import './App.css';

import Dashboard from './pages/Dashboard';
import DaemonDetail from './pages/DaemonDetail';
import VolumeManager from './pages/VolumeManager';
import Nav from './components/Header/Header';

function App(): React.ReactElement {
  return (
    <div className="App">
      <Router>
        <Routes>
          <Route path="/"                    element={<LandingPage />} />
          <Route path="/dashboard"           element={<Dashboard />} />
          <Route path="/daemon/:host"        element={<DaemonDetail />} />
          <Route path="/volume-directory"    element={<VolumeManager />} />
        </Routes>
      </Router>
    </div>
  );
}

const LandingPage: React.FC = (): React.ReactElement => {
  return (
    <div className="landing">
      <Nav />

      <div className="landing-hero">
        <div className="landing-badge">Docker Management Dashboard</div>
        <h1 className="landing-title">container-bay</h1>
        <p className="landing-subtitle">
          Remotely manage Docker containers, images, and volumes across multiple hosts — all from one place.
        </p>
        <div className="landing-actions">
          <Link className="landing-btn-primary" to="/dashboard">Open Dashboard</Link>
        </div>
      </div>

      <div className="landing-section">
        <h2 className="landing-section-title">What is Docker?</h2>
        <p className="landing-section-body">
          Docker is a platform for packaging, distributing, and running applications inside
          lightweight isolated environments called <strong>containers</strong>. A container bundles
          your application code together with all its dependencies (libraries, configs, runtime)
          so it runs identically on any machine — developer laptop, CI server, or production host.
        </p>
        <div className="landing-concepts">
          <div className="landing-concept-card">
            <div className="landing-concept-icon">▣</div>
            <div className="landing-concept-name">Container</div>
            <div className="landing-concept-desc">
              A running instance of an image. Isolated process with its own filesystem, network,
              and process space. Starts in milliseconds.
            </div>
          </div>
          <div className="landing-concept-card">
            <div className="landing-concept-icon">◈</div>
            <div className="landing-concept-name">Image</div>
            <div className="landing-concept-desc">
              A read-only template used to create containers. Built from a Dockerfile and stored
              in a registry. Composed of layered filesystems.
            </div>
          </div>
          <div className="landing-concept-card">
            <div className="landing-concept-icon">◉</div>
            <div className="landing-concept-name">Volume</div>
            <div className="landing-concept-desc">
              Persistent storage that lives outside the container filesystem. Volumes survive
              container removal and can be shared across containers.
            </div>
          </div>
        </div>
      </div>

      <div className="landing-section landing-section-alt">
        <h2 className="landing-section-title">What is a Docker Daemon?</h2>
        <p className="landing-section-body">
          The <strong>Docker daemon</strong> (<code>dockerd</code>) is the background service
          that does the heavy lifting — building images, running containers, managing networks
          and volumes. It exposes a REST API that clients (like this dashboard) talk to.
        </p>
        <div className="landing-daemon-diagram">
          <div className="landing-daemon-box client">
            <div className="landing-daemon-label">This Dashboard</div>
            <div className="landing-daemon-sub">REST API calls</div>
          </div>
          <div className="landing-daemon-arrow">
            <span className="landing-daemon-protocol">TCP (e.g. :2375)</span>
            <span className="landing-daemon-line">────────────→</span>
          </div>
          <div className="landing-daemon-box daemon">
            <div className="landing-daemon-label">dockerd</div>
            <div className="landing-daemon-sub">Remote host</div>
          </div>
        </div>
        <p className="landing-section-body">
          To allow remote access, the daemon must be started with <code>-H tcp://0.0.0.0:2375</code>
          (or a TLS-secured equivalent). Register the host address in the Dashboard to start managing it.
        </p>
      </div>

      <div className="landing-section">
        <h2 className="landing-section-title">How it works</h2>
        <div className="landing-features">
          <div className="landing-feature-card">
            <div className="landing-feature-icon">①</div>
            <div className="landing-feature-title">Register a Daemon</div>
            <ul className="landing-feature-list">
              <li>Go to Dashboard</li>
              <li>Enter the host address (IP:port)</li>
              <li>Give it an optional label</li>
            </ul>
            <Link className="landing-feature-cta" to="/dashboard">Go to Dashboard →</Link>
          </div>
          <div className="landing-feature-card">
            <div className="landing-feature-icon">②</div>
            <div className="landing-feature-title">Select a Daemon</div>
            <ul className="landing-feature-list">
              <li>Click any daemon in the list</li>
              <li>Opens the Containers tab by default</li>
              <li>Switch to Images or Logs tabs</li>
            </ul>
          </div>
          <div className="landing-feature-card">
            <div className="landing-feature-icon">③</div>
            <div className="landing-feature-title">Manage Containers</div>
            <ul className="landing-feature-list">
              <li>View all running/stopped containers</li>
              <li>Run, stop, restart, remove</li>
              <li>Tail logs, build images from Git</li>
            </ul>
          </div>
        </div>
      </div>

      <footer className="landing-footer">
        container-bay — Docker management over TCP
      </footer>
    </div>
  );
};

export default App;
