import * as React from 'react';


import { Actions } from 'actions/actions';

import YourLadders from 'components/YourLadders';
import NewLadder from 'components/NewLadder';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import API from 'api';

const mapStateToProps = (state: AppState) => ({
  ladders: state.ladders
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  createLadder: (name: string) => API.createLadder(name)
    .then(() => API.getLadders())
    .then(ladders => dispatch({
      type: Actions.SET_PLAYER_LADDERS,
      ladders
    }))
});

const Ladders = ({ 
  ladders, 
  createLadder 
}: { 
  ladders: LadderState, 
  createLadder: Function 
}) =>
  <div className="columns">
    <YourLadders owned={ladders.owned} playing={ladders.playing} className="column is-8" />
    <NewLadder createLadder={createLadder} className="column is-4" />
  </div>

export default connect(
  mapStateToProps, 
  mapDispatchToProps
)(Ladders)