import * as React from 'react';

import { connect } from 'react-redux';

import TitleBar from 'components/TitleBar';
import API from 'api';

type Props = {
  ladder: Ladder
}

const mapStateToProps = (state: AppState) => ({
  ladder: state.ladders.current
})

class LadderAdmin extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  render() {
    return (
      <div>
        <TitleBar ladder={this.props.ladder} />
      
        <div className="container">
          <h2 className="title is-3">Admin</h2>
          <div>Ladder season: {this.props.ladder.season}</div>

          <button className="button is-warning"
                  onClick={() => API.startNewSeason(this.props.ladder.id)}>
                  start new season
          </button>
        </div>
      </div>
    )
  }
}

export default connect(
  mapStateToProps
)(LadderAdmin)