import * as React from 'react';

import { Link } from 'react-router-dom';
import { firebase } from 'auth';

const ActiveLink = ({ href, children }: { href: string, children?: any }) => {
  const active = window.location.pathname === href

  return <li className={active ? 'is-active' : ''}>
    <Link to={href}>{children}</Link>
  </li>
}

const TitleBar = ({ ladder }: { ladder: Ladder }) => {
  const admin = firebase.signedIn() && ladder.ownerID === firebase.currentUser().uid;

  return (
    <section className="hero is-light is-small">
      <div className="hero-body">
        <div className="container">
          <h1 className="title is-4">{ladder.name}</h1>
          {ladder.season > 1 && <h2 className="subtitle is-7">season {ladder.season}</h2>}
        </div>
      </div>

      <div className="hero-foot is-hidden-mobile">
        <div className="container">
          <nav className="tabs is-boxed is-small">
            <ul>
              <ActiveLink href={`/ladder/${ladder.id}`}>table</ActiveLink>
              <ActiveLink href={`/chart/${ladder.id}`}>chart</ActiveLink>
              {admin && <ActiveLink href={`/admin/${ladder.id}`}>admin</ActiveLink>}
            </ul>
          </nav>
        </div>
      </div>
    </section>
  )
  
}

export default TitleBar;