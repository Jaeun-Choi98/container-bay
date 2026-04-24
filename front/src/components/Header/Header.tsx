import React from "react";
import { Link } from "react-router-dom";
import './Header.css'

function Nav(): React.ReactElement {
  return (
    <nav className="navbar">
      <span className="navbar-brand">
        <span className="navbar-brand-icon">▣</span>
        container-bay
      </span>
      <div className="navbar-links">
        <Link className="navbarMenu" to="/dashboard">Dashboard</Link>
        <Link className="navbarMenu" to="/volume-directory">Volume</Link>
      </div>
    </nav>
  );
}

export default Nav;
