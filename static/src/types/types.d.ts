/**
 * overall state for the app 
 * (built by combineReducers)
 */
type AppState = {
  user: UserState,
  view: ViewState
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

type LoginMode = 'register' | 'login';