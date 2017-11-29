import { AnyAction } from 'redux';
import { Actions } from 'actions/actions';

const initialState: LadderState = {
  owned: [],
  playing: []
}

export default function ladderReducer(state = initialState, action: AnyAction): LadderState {
  switch (action.type) {
    case Actions.SET_PLAYER_LADDERS: 
      return Object.assign({}, state, action.ladders);
    default:
      return state;
  }
}