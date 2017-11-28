import * as React from 'react';

import { connect } from 'react-redux';
import { Dispatch } from 'redux';

import Navbar from 'components/nav/Navbar';
import Login from 'components/Login';
import Ladders from 'components/Ladders';

import * as actions from 'actions/actions';
import { firebase } from 'auth';
import { User } from 'firebase';

const mapStateToProps = (state: AppState) => ({
  user: state.user,
  view: state.view
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLoginMode: (mode: LoginMode) => 
    dispatch(actions.setLoginMode(mode)),

  signIn: (email: string, password: string): Promise<any> => 
    firebase.signIn(email, password),

  signOut: () => firebase.signOut().then(() => dispatch({
    type: actions.SIGN_OUT
  })),

  register: (username: string, email: string, password: string): Promise<any> =>
    firebase.register(username, email, password)
      .then((user: User) => dispatch({
        type: actions.SIGN_IN,
        username: user.displayName
      })),

  userFormInput: (field: string, value: string) => dispatch({
    type: actions.USER_FORM_INPUT,
    field,
    value
  })
});

// type of the props passed into Main (as built by `connect`)
type MainProps = { 
  user: UserState, 
  view: ViewState,
  setLoginMode: (mode: LoginMode) => any,
  register: (username: string, email: string, password: string) => Promise<any>,
  signIn: (email: string, password: string) => Promise<any>,
  signOut: () => Promise<any>,
  userFormInput: (field: string, value: string) => void
};

const Main = ({ user, view, setLoginMode, register, signIn, signOut, userFormInput }: MainProps) => {
  return (
    <div>
      <Navbar username={user.username} signOut={signOut} />

      <section className="section">
        <div className="container">

        {user.signedIn ? 
          <Ladders /> :

          <Login 
            input={user.formInput}
            mode={view.loginMode}
            selectLogin={() => setLoginMode('login')}
            selectRegister={() => setLoginMode('register')}
            register={register}
            signIn={signIn} 
            inputName={(username: string) => userFormInput('username', username)} 
            inputEmail={(email: string) => userFormInput('email', email)} 
            inputPassword={(password: string) => userFormInput('password', password)} />  
        }
        </div>
      </section>
    </div>
  );
};

export default connect(
  mapStateToProps, 
  mapDispatchToProps
)(Main)