export enum Actions {
  SET_LOGIN_MODE = 'SET_LOGIN_MODE',
  SIGN_IN = 'SIGN_IN',
  SIGN_OUT = 'SIGN_OUT',
  USER_FORM_INPUT = 'USER_FORM_INPUT',

  SET_PLAYER_LADDERS = 'SET_PLAYER_LADDERS'
}

export function signIn(username: string) {
  return {
    type: Actions.SIGN_IN,
    username: username
  }
}

export function setLoginMode(mode: LoginMode) {
  return {
    type: Actions.SET_LOGIN_MODE,
    mode: mode
  }
}