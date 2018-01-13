import { AnyAction } from 'redux';
import { Actions } from 'actions/actions';

const initialState: ModalState = {
  active: false,
  level: 'info',
  message: ''
}

export default function ladderReducer(state = initialState, action: AnyAction): ModalState {
  switch (action.type) {
    case Actions.SHOW_MODAL: 
      return Object.assign({}, state, {
        active: true,
        level: action.level,
        message: action.message
      });

    case Actions.HIDE_MODAL:
      return Object.assign({}, state, {
        active: false,
        message: ''
      });

    default:
      return state;
  }
}