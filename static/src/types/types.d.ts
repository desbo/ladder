// API URL, set in webpack config
declare const API_URL: string;

type Ladder = {
  name: string,
  key: string,
  created: string,
  players: Array<{
    position: number,
    name: string,
    wins: number,
    losses: number
  }>
}

/**
 * overall state for the app 
 * (built by combineReducers)
 */
type AppState = {
  user: UserState,
  view: ViewState,
  ladders: LadderState
}

// Login/Registration input
type LoginFormInput = {
  username?: string,
  email: string,
  password: string
}

type UserState = {
  formInput: LoginFormInput,
  signedIn: boolean,
  username: null | string
}

type ViewState = {
  loginMode: LoginMode
}

type LadderState = {
  owned: Array<Ladder>,
  playing: Array<Ladder>
}

type LoginMode = 'register' | 'login';