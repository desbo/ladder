import * as React from 'react';
import { BrowserRouter, Route, Redirect, Switch } from 'react-router-dom';
import { connect, Dispatch } from 'react-redux';

import { firebase } from 'auth';

import { Actions } from 'actions/actions';

import Main from 'components/Main';
import Navbar from 'components/nav/Navbar';
import ViewLadder from 'components/ViewLadder';
import Modal from 'components/util/Modal';

const mapStateToProps = (state: AppState) => ({
  username: state.user.username,
  modal: state.modal
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  signOut: () => firebase.signOut().then(() => dispatch({ 
    type: Actions.SIGN_OUT
  })),

  closeModal: () => dispatch({ type: Actions.HIDE_MODAL })
});

const AuthedRoute = ({ username, ...props }: { username: string, [prop: string]: any }) => {
  if (!username) 
    return <Redirect to="/"/>
  else 
    return <Route {...props} />
}

const Router = ({ username, modal, signOut, closeModal }: { username: string, modal: ModalState, signOut: Function, closeModal: Function }) => {
  return (<BrowserRouter>
    <div>
      {modal.active && (
        <Modal level={modal.level} active={modal.active} close={closeModal}>
          {modal.message}
        </Modal>
      )}

      <Navbar username={username} signOut={signOut} />

      <section className="section">
        <div className="container">
          <Switch>
            <Route exact path="/" component={Main} />
            <AuthedRoute username={username} path="/ladder/:id" component={ViewLadder} />
          </Switch>
        </div>
      </section>
    </div>
    </BrowserRouter>);
}
  

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Router)