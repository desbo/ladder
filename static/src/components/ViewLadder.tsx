import * as React from 'react';

import { RouteProps, match, Redirect } from 'react-router';
import { ReactNode } from 'react';

import { Actions, setCurrentLadder, selectOpponent, clearOpponent } from 'actions/actions';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import Table from 'components/ladder/Table'
import SubmitGame from 'components/ladder/SubmitGame'
import PlayerDropdown from 'components/ladder/PlayerDropdown'
import Login from 'components/Login';

import API from 'api';

const mapStateToProps = (state: AppState) => ({
  ladder: state.ladders.current,
  opponent: state.ladders.opponent,
  user: state.user
})

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLadder: (ladder: Ladder) => dispatch(setCurrentLadder(ladder)),
  selectOpponent: (opponent: LadderPlayer) => dispatch(selectOpponent(opponent)),
  clearOpponent: () => dispatch(clearOpponent())
})

type Props = {
  ladder: Ladder,
  opponent: LadderPlayer,
  user: UserState,
  setLadder: (ladder: Ladder) => any,
  selectOpponent: (opponent: LadderPlayer) => any,
  clearOpponent: () => any,
  match: match<{ id: string }>
};

type State = {
  failed: boolean
}

class ViewLadder extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      failed: false
    };
  }

  componentDidMount() {
    this.fetch()
  }

  fetch() {
    return API.getLadder(this.props.match.params.id)
      .then(ladder => this.props.setLadder(ladder))
      .catch(() => this.setState({ failed: true }))
  }

  join() {
    return API.joinLadder(this.props.ladder.id)
      .then(this.fetch.bind(this))
  }

  render() {
    if (this.props.ladder) {
      return (
        <div>
          <section className="hero is-light is-small">
            <div className="hero-body">
              <div className="container">
                <h1 className="subtitle is-4">{this.props.ladder.name}</h1>
              </div>
            </div>
          </section>

          <section className="section">
            <div className="container">
              <div className="columns reverse-row-order">
                {this.props.user.signedIn && 
                  <div className="column is-5">
                    <div style={{"marginBottom": "-8px"}} className="columns level is-mobile">
                      <div className="column is-6">
                        <h2 className="subtitle is-4 level-item level-left">submit game</h2>
                      </div>
                      <div className="column is-6">
                        <PlayerDropdown 
                          players={this.props.ladder.players.filter(p => p.name !== this.props.user.username )}
                          onSelect={this.props.selectOpponent} />
                      </div>
                    </div>
                    
                    <SubmitGame 
                      user={this.props.user} 
                      ladder={this.props.ladder} 
                      opponent={this.props.opponent}
                      onSubmit={() => this.fetch().then(this.props.clearOpponent)} />
                  </div>
                }

                {!this.props.user.signedIn && 
                  <div className="column is-5">
                    <h2 className="subtitle is-4">join this ladder</h2>
                    <Login fullWidth={true} registerOnly={true}
                          onRegister={this.join.bind(this)} />
                  </div>
                }

                <div className="column is-7">
                  <Table ladder={this.props.ladder} />
                </div>
              </div>
            </div>
          </section>
        </div>
      );
    } else if (this.state.failed) {
      return <Redirect to="/" />
    } else {
      return <div>getting ladder...</div>
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ViewLadder)