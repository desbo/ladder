import * as React from 'react';

import { Actions } from 'actions/actions';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import API from 'api';
import { RouteProps, match } from 'react-router';

const mapStateToProps = (state: AppState) => ({
  ladders: state.ladders.owned.concat(state.ladders.playing)
})

const ViewLadder = ({ 
  match, 
  ladders,
}: { 
  match: match<{ id: string }>,
  ladders: Array<Ladder>
}) => {
  const ladderID = match.params.id;
  const ladder = ladders.find(l => l.key === ladderID);

  return (
    <div>
      {ladder && ladder.name}
    </div>
  );  
}

export default connect(
  mapStateToProps
)(ViewLadder)