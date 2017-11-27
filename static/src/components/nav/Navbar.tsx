import * as React from 'react';

import UserMenu from 'components/nav/UserMenu';

export default function Navbar({
  username,
  signOut
}: { 
  username: null | string,
  signOut: Function
}) {
  return (
    <nav className="navbar" role="navigation" aria-label="main navigation">
      <div className="container">      
        <div className="navbar-brand">
          <div className="navbar-item title is-4">table tennis ladder 🏓</div>
        </div>

        {username && <UserMenu username={username} signOut={signOut} /> }
      </div>
    </nav>
  );
}