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
    console.log('registering');
    console.log(username, email, password);
    return firebase.register(username, email, password).then(x => console.log(x));
  }
});

// type of the props passed into Main (as built by `connect`)
type MainProps = { 
  user: UserState, 
  view: ViewState,
  setLoginMode: (mode: LoginMode) => any,
  register: (username: string, email: string, password: string) => Promise<any>,
  signIn: (email: string, password: string) => Promise<any>
};

const Main = ({ user, view, setLoginMode, register, signIn }: MainProps) => {
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
            signIn={signIn} />  
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