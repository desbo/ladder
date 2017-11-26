import { SET_LOGIN_MODE } from 'actions/actions';
import { AnyAction } from 'redux';

const initialState: ViewState = {
  loginMode: 'login'
};

export default function viewReducer(state = initialState, action: AnyAction): ViewState {
  switch (action.type) {
    case SET_LOGIN_MODE:
      return Object.assign({}, state, {
        loginMode: action.mode
      });

    default: 
      return state;
  }
}