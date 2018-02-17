import { AnyAction } from 'redux';
import { Actions } from 'actions/actions';

const initialState: Maybe<ChartData> = undefined

export default function chartReducer(state = initialState, action: AnyAction): Maybe<ChartData> {
  switch (action.type) {
    case Actions.SET_CHART_DATA:
      return action.data
    default:
      return null
  }
}