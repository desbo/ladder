import * as React from 'react';

export default function Navbar({
  username
}: { username: null | string }) {
  return (
    <nav className="navbar" role="navigation" aria-label="main navigation">
      <div className="container">      
        <div className="navbar-brand">
          <div className="navbar-item title is-4">table tennis ladder 🏓</div>
        </div>

        {username &&
          <div className="navbar-right">{username}</div> 
        }
      </div>
    </nav>
  );
}