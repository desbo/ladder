import * as React from 'react';

import { connect } from 'react-redux';
import { Dispatch } from 'redux';

import Login from 'components/Login';
import Ladders from 'components/Ladders';

const mapStateToProps = (state: AppState) => ({
  user: state.user
});

type MainProps = { 
  user: UserState
};

const Main = ({ user }: MainProps) => {
  return (
    <section className="section">
      <div className="container">
        {user.signedIn ? 
          <Ladders /> :
          <Login />  
        }
      </div>
    </section>
  );
};

export default connect(mapStateToProps)(Main);