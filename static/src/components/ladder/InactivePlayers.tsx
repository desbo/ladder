import * as React from 'react';

const InactivePlayers = ({ ladder }: { ladder: Ladder }) => 
  <div>
    Inactive players: {
      ladder.players.filter(p => !p.active).map(p => p.name).join(', ')
    }
  </div>;
  
export default InactivePlayers;