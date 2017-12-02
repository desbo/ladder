import { AnyAction } from 'redux';
import { Actions } from 'actions/actions';

const initialState: LadderState = {
  owned: [],
  playing: []
}

const sortLadder = (a: Ladder, b: Ladder) => 
  Date.parse(b.created) - Date.parse(a.created)

export default function ladderReducer(state = initialState, action: AnyAction): LadderState {
  switch (action.type) {
    case Actions.SET_PLAYER_LADDERS: 
      return Object.assign({}, state, {
        owned: action.ladders.owned.sort(sortLadder),
        playing: action.ladders.playing.sort(sortLadder)
      });

    default:
      return state;
  }
}