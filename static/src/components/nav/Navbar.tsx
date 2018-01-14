import * as React from 'react';

import UserMenu from 'components/nav/UserMenu';
import { Link } from 'react-router-dom';

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
          <div className="navbar-item title is-4"><Link to="/">🏓</Link></div>
        </div>

        {username && <UserMenu username={username} signOut={signOut} /> }
      </div>
    </nav>
  );
}