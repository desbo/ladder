import * as React from 'react';

import InputField from 'components/form/InputField';

function format(date: string): string {
  return new Date(Date.parse(date)).toLocaleString();
}

const YourLadders = ({ 
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
        {owned.map(ladder => {

          return <tr key={ladder.key}>
            <td><a href={ladder.key}>{ladder.name}</a></td>
            <td key={`${ladder.key}-date`}>{format(ladder.created)}</td>
          </tr>
        })}
      </tbody>
    </table>
  </div>
);


export default YourLadders;