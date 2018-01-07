import * as React from 'react';
import { BrowserRouter, Route, Redirect } from 'react-router-dom';
import { connect, Dispatch } from 'react-redux';

import { firebase } from 'auth';

import { Actions } from 'actions/actions';

import Main from 'components/Main';
import Navbar from 'components/nav/Navbar';
import ViewLadder from 'components/ViewLadder';


const mapStateToProps = (state: AppState) => ({
  username: state.user.username
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  signOut: () => firebase.signOut().then(() => dispatch({
    type: Actions.SIGN_OUT
  })),
});

const AuthedRoute = ({ username, ...props }: { username: string, [prop: string]: any }) => {
  if (!username) 
    return <Redirect to="/"/>
  else 
    return <Route {...props} />
}

const Router = ({ username, signOut }: { username: string, signOut: Function }) =>
  <BrowserRouter>
    <div>
      <Navbar username={username} signOut={signOut} />

      <section className="section">
        <div className="container">
          <Route exact path="/" component={Main} />
          <AuthedRoute username={username} path="/ladder/:id" component={ViewLadder} />
        </div>
      </section>
    </div>
  </BrowserRouter>;

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Router)