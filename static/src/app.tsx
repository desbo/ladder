import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Store, createStore, combineReducers } from 'redux';
import { User } from 'firebase';

import * as auth from 'auth';

import { signIn } from 'actions/actions';
import Main from 'components/Main';

import userReducer from 'reducers/user';
import viewReducer from 'reducers/view';

const store: Store<AppState> = createStore(combineReducers({
  user: userReducer,
  view: viewReducer,
}));

const firebaseApp = auth.initFirebase();

firebaseApp.app.auth().onAuthStateChanged((user: null | User) => {
  if (user) store.dispatch(signIn(user.displayName))
});

ReactDOM.render(
  <Provider store={store}>
    <Main />
  </Provider>,
  document.getElementById('app')
)