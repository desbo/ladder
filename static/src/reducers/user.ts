import { AnyAction } from 'redux';
import { SIGN_IN, SIGN_OUT, USER_FORM_INPUT } from 'actions/actions';

const initialState: UserState = {
  formInput: {
    username: '',
    email: '',
    password: ''
  },
  signedIn: false,
  username: null
}

export default function userReducer(state = initialState, action: AnyAction): UserState {
  switch (action.type) {
    case USER_FORM_INPUT: 
      return Object.assign({}, state, {
        formInput: Object.assign({}, state.formInput, { [action.field]: action.value })
      });

    case SIGN_IN: 
      return Object.assign({}, state, {
        signedIn: true,
        username: action.username
      });

    case SIGN_OUT:
      return Object.assign({}, state, {
        signedIn: false,
        username: null
      });

    default:
      return state;
  }
}