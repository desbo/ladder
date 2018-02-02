// constants set by webpack define plugin
declare const API_URL: string;
declare const FIREBASE_CONFIG: Object;


type LadderPlayer = {
  key: string,
  position: number,
  name: string,
  wins: number,
  losses: number,
  rating: number,
  winRate: number,
  active: boolean
}

type Ladder = {
  name: string,
  id: string,
  created: string,
  players: Array<LadderPlayer>
}

type Player = {
  name: string,
  rating: number
}

type PlayerResult = {
  player: Player,
  score: number
}

type Game = {
  id: string,
  date: string,
  player1: PlayerResult,
  player2: PlayerResult
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
  inlineFormInput: LoginFormInput,
  signedIn: boolean,
  username: null | string
}

type ViewState = {
  loginMode: LoginMode
}

type LadderState = {
  owned: Array<Ladder>,
  playing: Array<Ladder>,
  current: Ladder,
  opponent: LadderPlayer,
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