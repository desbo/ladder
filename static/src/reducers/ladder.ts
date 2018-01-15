import { AnyAction } from 'redux';
import { Actions } from 'actions/actions';
import { access } from 'fs';

const initialState: LadderState = {
  owned: [],
  playing: [],
  current: undefined,
  opponent: undefined
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

    case Actions.SET_CURRENT_LADDER: 
      return Object.assign({}, state, {
        current: action.ladder
      });

    case Actions.SELECT_OPPONENT:
      return Object.assign({}, state, {
        opponent: action.opponent
      });

    default:
      return state;
  }
}