// API URL, set in webpack config
declare const API_URL: string;

type LadderPlayer = {
  position: number,
  name: string,
  wins: number,
  losses: number,
  rating: number
}

type Ladder = {
  name: string,
  id: string,
  created: string,
  players: Array<LadderPlayer>
}

/**
 * overall state for the app 
 * (built by combineReducers)
 */
type AppState = {
  user: UserState,
  view: ViewState,
  ladders: LadderState,
  modal: ModalState
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

type ModalMessageLevel = 'info' | 'error';

type ModalState = {
  level: ModalMessageLevel,
  message: string,
  active: boolean
}

type LoginMode = 'register' | 'login';

type AppError = {
  message: string
}