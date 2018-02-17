import * as React from 'react';

import { match, Redirect } from 'react-router';
import { ReactNode } from 'react';

import { Actions, setCurrentLadder, selectOpponent, clearOpponent } from 'actions/actions';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import TitleBar from 'components/TitleBar'
import Table from 'components/ladder/Table'
import SubmitGame from 'components/ladder/SubmitGame'
import PlayerDropdown from 'components/ladder/PlayerDropdown'
import Login from 'components/Login';

import API from 'api';
import InactivePlayers from 'components/ladder/InactivePlayers';

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
  failed: boolean,
  joining: boolean
}

class ViewLadder extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      failed: false,
      joining: false,
    };
  }

  componentDidMount() {
    if (!this.props.ladder) {
      this.fetch()
    }
  }

  fetch() {
    return API.getLadder(this.props.match.params.id)
      .then(ladder => this.props.setLadder(ladder))
      .catch(() => this.setState({ failed: true }))
  }

  join() {
    this.setState({ joining: true });

    return API.joinLadder(this.props.ladder.id)
      .then(this.fetch.bind(this))
      .then(() => this.setState({ joining: false }));
  }

  render() {
    if (this.props.ladder) {
      return (
        <div>
          <TitleBar ladder={this.props.ladder} />

          <section className="section">
            <div className="container">
              <div className="columns reverse-row-order">
                {this.props.user.signedIn && 
                    (this.props.ladder.players.some(p => p.name == this.props.user.username) ?
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
                      </div> :
                        
                      <div className="column is-5 has-text-centered">
                        <button 
                          onClick={() => API.joinLadder(this.props.ladder.id)
                            .then(() => API.getLadder(this.props.ladder.id))
                            .then(ladder => this.props.setLadder(ladder))}
                          className={`button is-primary ${this.state.joining ? 'is-loading' : ''}`}
                          disabled={this.state.joining}>
                            join this ladder
                        </button>
                      </div>
                    )
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

                  {this.props.ladder.players.filter(p => !p.active).length > 0 &&
                    <InactivePlayers ladder={this.props.ladder} />
                  }
                </div>
              </div>
            </div>
          </section>
        </div>
      );
    } else if (this.state.failed) {
      return <Redirect to="/" />
    } else {
      return <div className="container">getting ladder...</div>
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ViewLadder)