import * as React from 'react';

import InputField from 'components/form/InputField';
import { Link } from 'react-router-dom';

function format(date: string): string {
  return new Date(Date.parse(date)).toLocaleString();
}

function distinct(ls: Array<Ladder>): Array<Ladder> {
  return ls.reduce((ladders, ladder) => {
    if (ladders.some(l => l.id === ladder.id))
      return ladders
    else 
      return ladders.concat(ladder)
  }, [])
}

const LadderList = ({ 
  owned, 
  playing, 
  className 
}: { 
  owned: Array<Ladder>,
  playing: Array<Ladder>,
  className: string 
}) => (
  <div className={className}>
    <h2 className="title is-4">your ladders</h2>
    <table className="table is-hoverable is-fullwidth">
      <thead>
        <tr>
          <th style={ { width: "70%" } }>name</th>
          <th>created</th>
        </tr>
      </thead>

      <tbody>
        {distinct(owned.concat(playing)).map(ladder => {
          return <tr key={ladder.id}>
            <td><Link to={`/ladder/${ladder.id}`}>{ladder.name}</Link></td>
            <td key={`${ladder.id}-date`}>{format(ladder.created)}</td>
          </tr>
        })}
      </tbody>
    </table>
  </div>
);

export default LadderList;
