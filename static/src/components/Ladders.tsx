import * as React from 'react';

import YourLadders from 'components/YourLadders';
import NewLadder from 'components/NewLadder';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import { local as localAPI } from 'api';

const mapStateToProps = (state: AppState) => ({
  user: state.user
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  createLadder: (name: string) => localAPI.createLadder(name)
    .then(r => console.log(r))
});

const Ladders = ({ 
  user, 
  createLadder 
}: { 
  user: UserState, 
  createLadder: Function 
}) =>
  <div className="columns">
    <YourLadders className="column" />
    <NewLadder createLadder={createLadder} className="column" />
  </div>

export default connect(
  mapStateToProps, 
  mapDispatchToProps
)(Ladders)