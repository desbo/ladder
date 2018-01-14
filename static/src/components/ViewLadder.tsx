import * as React from 'react';

import { Actions, setCurrentLadder } from 'actions/actions';

import { Dispatch } from 'redux';
import { connect } from 'react-redux';

import Table from 'components/ladder/Table'

import API from 'api';
import { RouteProps, match, Redirect } from 'react-router';
import { ReactNode } from 'react';
import Login from 'components/Login';

type Props = {
  ladder: Ladder,
  signedIn: boolean,
  setLadder: (ladder: Ladder) => any,
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
              <div className="columns">
                <div className="column is-6">
                  <Table ladder={this.props.ladder} />
                </div>

                <div className="column is-offset-2 is-4">
                {this.props.signedIn && 
                  <h2 className="subtitle is-4">submit game</h2>
                }

                {!this.props.signedIn && 
                <div>
                  <h2 className="subtitle is-4">join this ladder</h2>
                  <Login fullWidth={true} registerOnly={true}
                         onRegister={this.join.bind(this)} />
                </div>
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
      return <div>getting ladder...</div>
    }
  }
}

const mapStateToProps = (state: AppState) => ({
  ladder: state.ladders.current,
  signedIn: state.user.signedIn
})

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLadder: (ladder: Ladder) => dispatch(setCurrentLadder(ladder))
})

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ViewLadder)