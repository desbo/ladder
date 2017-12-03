import { Actions } from 'actions/actions';
import { AnyAction } from 'redux';

const initialState: ViewState = {
  loginMode: 'login'
};

export default function viewReducer(state = initialState, action: AnyAction): ViewState {
  switch (action.type) {
    case Actions.SET_LOGIN_MODE:
      return Object.assign({}, state, {
        loginMode: action.mode
      });

    default: 
      return state;
  }
}