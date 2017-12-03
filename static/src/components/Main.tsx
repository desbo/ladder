import * as React from 'react';

import { connect } from 'react-redux';
import { Dispatch } from 'redux';

import Login from 'components/Login';
import Ladders from 'components/Ladders';

import { Actions , setLoginMode } from 'actions/actions';
import { firebase } from 'auth';
import { User } from 'firebase';

const mapStateToProps = (state: AppState) => ({
  user: state.user,
  view: state.view
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLoginMode: (mode: LoginMode) => 
    dispatch(setLoginMode(mode)),

  signIn: (email: string, password: string): Promise<any> => 
    firebase.signIn(email, password),

  register: (username: string, email: string, password: string): Promise<any> =>
    firebase.register(username, email, password)
      .then((user: User) => dispatch({
        type: Actions.SIGN_IN,
        username: user.displayName
      })),

  userFormInput: (field: string, value: string) => dispatch({
    type: Actions.USER_FORM_INPUT,
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
  userFormInput: (field: string, value: string) => void
};

const Main = ({ user, view, setLoginMode, register, signIn, userFormInput }: MainProps) => {
  return (
    <div>
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
  );
};

export default connect(
  mapStateToProps, 
  mapDispatchToProps
)(Main)