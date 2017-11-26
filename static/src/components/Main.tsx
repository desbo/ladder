import * as React from 'react';

import { connect } from 'react-redux';
import { Dispatch } from 'redux';

import Navbar from 'components/Navbar';
import Login from 'components/Login';

import * as actions from 'actions/actions';
import { firebase } from 'auth';

const mapStateToProps = (state: AppState) => ({
  user: state.user,
  view: state.view
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLoginMode: (mode: LoginMode) => dispatch(actions.setLoginMode(mode)),

  signIn: (email: string, password: string): Promise<any> => {
    return firebase.signIn(email, password).then(x => console.log(x));
  },

  register: (username: string, email: string, password: string): Promise<any> => {
    return firebase.register(username, email, password).then(x => console.log(x));
  },

  userFormInput: (field: string, value: string): void => {
    dispatch({
      type: actions.USER_FORM_INPUT,
      field,
      value
    });
  }
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
      <Navbar />

      <section className="section">
        <div className="container">

        {user.signedIn ? 
          <h1>Hello user!</h1> :

          <Login 
            input={user.formInput}
            mode={view.loginMode}
            selectLogin={() => setLoginMode('login')}
            selectRegister={() => setLoginMode('register')}
            register={register}
            signIn={signIn} 
            inputName={(username: string) => userFormInput('username', name)} 
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