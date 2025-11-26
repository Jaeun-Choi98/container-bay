import React from 'react';
import JSX from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import logo from './logo.svg';
import './App.css';

import Home from './Home'


function App(): JSX.ReactElement {
  return (
    <div className="App">
      <Router>
        <Routes>
          <Route path="/" element={<LearnReact />}></Route>
          <Route path="/home" element={<Home />}></Route>
        </Routes>
      </Router>

    </div>
  );
}

const LearnReact: React.FC = () => {
  return (
    <div>
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Start <code>src/App.tsx</code> and save to reload.
        </p>
        <Link to="/home"
          className="App-link"
          target="_blank"

        >
          Go Container Bay
        </Link>
      </header>
    </div>
  )
}

export default App;
