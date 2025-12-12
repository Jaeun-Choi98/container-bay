import React from "react";
import { Link } from "react-router-dom";
import '../css/Nav.css'

function Nav(): React.ReactElement {

  return (
    <div>
      <div className="navbar">
        <Link className="navbarMenu" to={'/container'}>Container</Link>
        <Link className="navbarMenu" to={'/image'}>Image</Link>
        <Link className="navbarMenu" to={'/volume-directory'}>Volume Directory</Link>
      </div>
    </div>
  );
}

export default Nav;