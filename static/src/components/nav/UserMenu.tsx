import * as React from 'react';
import { MouseEvent } from 'react';

const UserMenu = ({ username, signOut }: { username: string, signOut: Function }) => (
  <div className="navbar-menu">
  
    <div className="navbar-end">
      <div className="navbar-item has-dropdown is-hoverable">
        <div className="navbar-link">{username}</div>
        <div className="navbar-dropdown">
          <a className="navbar-item" onClick={() => signOut()}>Sign out</a>
        </div>
      </div>
    </div>
  </div>
);

export default UserMenu;