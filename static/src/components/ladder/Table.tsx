import * as React from 'react';

const RatingChange = ({change}: {change:number}) => {
  const up = change > 0;
  return (
    <span className={`rating-change has-text-right is-size-7 ${up ? 'has-text-success' : 'has-text-danger'}`}>
      <span className="icon">
        <i className={`fas ${up ? 'fa-level-up-alt' : 'fa-level-down-alt'}`}></i>
      </span>
      {change.toString().replace('-', '')}
    </span>
  )
}
  

const Table = ({ ladder }: { ladder: Ladder }) => 
  <table className="table is-striped is-hoverable is-fullwidth">
    <thead>
      <tr>
        <th><abbr title="position"></abbr></th>
        <th>Name</th>
        <th>Played</th>
        <th className="is-hidden-mobile">Won</th>
        <th className="is-hidden-mobile">Lost</th>
        <th>Win %</th>
        <th>Rating</th>
      </tr>
    </thead>

    <tbody>
      {ladder.players.filter(p => p.active).map(p => <tr key={p.name}>
        <td>{p.position}</td>
        <td>{p.name}</td>
        <td>{p.wins + p.losses}</td>
        <td className="is-hidden-mobile">{p.wins}</td>
        <td className="is-hidden-mobile">{p.losses}</td>
        <td>{Math.round(p.winRate * 100)}%</td>
        <td>{p.rating} {p.lastRatingChange !== 0 && <RatingChange change={p.lastRatingChange} />}</td>
      </tr>)}
    </tbody>
  </table>

export default Table;