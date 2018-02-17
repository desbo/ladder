import * as React from 'react';

import { Link } from 'react-router-dom';

const ActiveLink = ({ href, children }: { href: string, children?: any }) => {
  const active = window.location.pathname === href

  return <li className={active ? 'is-active' : ''}>
    <Link to={href}>{children}</Link>
  </li>
}

const TitleBar = ({ ladder }: { ladder: Ladder }) => {
  return (
    <section className="hero is-light is-small">
      <div className="hero-body">
        <div className="container">
          <h1 className="subtitle is-4">{ladder.name}</h1>
        </div>
      </div>

      <div className="hero-foot">
        <div className="container">
          <nav className="tabs is-boxed is-small">
            <ul>
              <ActiveLink href={`/ladder/${ladder.id}`}>table</ActiveLink>
              <ActiveLink href={`/chart/${ladder.id}`}>chart</ActiveLink>
            </ul>
          </nav>
        </div>
      </div>
    </section>
  )
  
}

export default TitleBar;