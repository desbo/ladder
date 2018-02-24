import * as React from 'react';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';
import { match } from 'react-router';

import { setCurrentLadder } from 'actions/actions';

import TitleBar from 'components/TitleBar';
import API from 'api';

type Props = {
  ladder: Ladder,
  setLadder: (ladder: Ladder) => any,
  match: match<{ id: string }>
}

const mapStateToProps = (state: AppState) => ({
  ladder: state.ladders.current
})

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLadder: (ladder: Ladder) => dispatch(setCurrentLadder(ladder))
})

class LadderAdmin extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  componentDidMount() {
    this.fetch();
  }

  fetch() {
    return API.getLadder(this.props.match.params.id)
      .then(ladder => this.props.setLadder(ladder));
  }

  render() {
    if (this.props.ladder) return (
      <div>
        <TitleBar ladder={this.props.ladder} />
      
        <div className="container">
          <h2 className="title is-3">Admin</h2>
          <div>Ladder season: {this.props.ladder.season}</div>

          <button className="button is-warning"
                  onClick={() => 
                    API.startNewSeason(this.props.ladder.id)
                      .then(() => this.fetch())
          }>
                  start new season
          </button>
        </div>
      </div>
    )
    else return (
      <div className="container">getting ladder...</div>
    )
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(LadderAdmin)