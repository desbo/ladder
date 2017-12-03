import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Store, createStore, combineReducers } from 'redux';
import { User } from 'firebase';

import * as auth from 'auth';
import API from 'api';

import { Actions, signIn, setPlayerLadders } from 'actions/actions';

import Router  from 'components/router';

import userReducer from 'reducers/user';
import viewReducer from 'reducers/view';
import ladderReducer from 'reducers/ladder';

const firebaseApp = auth.initFirebase();

const devToolEnhancer: any = 
  (window as any)['__REDUX_DEVTOOLS_EXTENSION__'] && 
  (window as any)['__REDUX_DEVTOOLS_EXTENSION__']();

const store: Store<AppState> = createStore(
  combineReducers({
    user: userReducer,
    view: viewReducer,
    ladders: ladderReducer
  }),
  devToolEnhancer
);

const renderMain = () => ReactDOM.render(
  <Provider store={store}>
    <Router />
  </Provider>,

  document.getElementById('app')
);

firebaseApp.app.auth().onAuthStateChanged((user: null | User) => {
  if (user) {
    store.dispatch(signIn(user.displayName));

    return API.getLadders()
      .then(ladders => store.dispatch(setPlayerLadders(ladders)))
      .then(renderMain);
  }

  renderMain();
});