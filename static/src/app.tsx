import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Store, createStore, combineReducers } from 'redux';
import { User } from 'firebase';

import * as auth from 'auth';
import API from 'api';

import { Actions, signIn } from 'actions/actions';
import Main from 'components/Main';

import userReducer from 'reducers/user';
import viewReducer from 'reducers/view';
import ladderReducer from 'reducers/ladder';

const store: Store<AppState> = createStore(combineReducers({
  user: userReducer,
  view: viewReducer,
  ladders: ladderReducer
}));

const firebaseApp = auth.initFirebase();

firebaseApp.app.auth().onAuthStateChanged((user: null | User) => {
  if (user) {
    store.dispatch(signIn(user.displayName));
    API.getLadders().then(ladders => store.dispatch({
      type: Actions.SET_PLAYER_LADDERS,
      ladders
    }));
  }

  ReactDOM.render(
    <Provider store={store}>
      <Main />
    </Provider>,
    document.getElementById('app')
  )
});