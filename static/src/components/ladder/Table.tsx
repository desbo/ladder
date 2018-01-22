import * as React from 'react';

const Table = ({ ladder }: { ladder: Ladder }) => 
  <table className="table is-striped is-hoverable is-fullwidth">
    <thead>
      <tr>
        <th><abbr title="position"></abbr></th>
        <th>Name</th>
        <th>Played</th>
        <th>Won</th>
        <th>Lost</th>
        <th>Win %</th>
        <th>Rating</th>
      </tr>
    </thead>

    <tbody>
      {ladder.players.map(p => <tr key={p.name}>
        <td>{p.position}</td>
        <td>{p.name}</td>
        <td>{p.wins + p.losses}</td>
        <td>{p.wins}</td>
        <td>{p.losses}</td>
        <td>{Math.round(p.winRate * 100)}%</td>
        <td>{p.rating}</td>
      </tr>)}
    </tbody>
  </table>

export default Table;