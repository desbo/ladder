import * as React from 'react';

import { Actions } from 'actions/actions';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import API from 'api';
import { RouteProps, match } from 'react-router';

const mapStateToProps = (state: AppState) => ({
  ladders: state.ladders.owned.concat(state.ladders.playing)
})

const signupURL = (ladder: Ladder): string => {
  const link = document.createElement("a")
  link.setAttribute('href', `/join/${ladder.id}`)
  return link.href;
}

const copy = (text: string) => {
  const el = document.createElement('textarea');
  el.innerHTML = text
  el.select()
  document.execCommand('copy');
}

const ViewLadder = ({ 
  match, 
  ladders,
}: { 
  match: match<{ id: string }>,
  ladders: Array<Ladder>
}) => {
  const ladder: Ladder | undefined = ladders.find(l => l.id === match.params.id);

  if (ladder) return (
    <div className="columns">
      <div className="column is-6">
        <h1 className="title">{ladder.name}</h1>
        <table className="table is-striped is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th><abbr title="position"></abbr></th>
              <th>Name</th>
              <th>Played</th>
              <th>Won</th>
              <th>Lost</th>
            </tr>
          </thead>

          <tbody>
            {ladder.players.map(p => <tr key={p.name}>
              <td>{p.position}</td>
              <td>{p.name}</td>
              <td>{p.wins + p.losses}</td>
              <td>{p.wins}</td>
              <td>{p.losses}</td>
            </tr>)}
          </tbody>
        </table>
      </div>

      <div className="column is-4 is-offset-2">
        <h2 className="title is-4">signup link</h2>
        <span className="tag is-light">{signupURL(ladder)}</span>
        <a onClick={copy.bind(null, signupURL(ladder))} className="button is-small">copy</a>
      </div>
    </div>
  );
}

export default connect(
  mapStateToProps
)(ViewLadder)