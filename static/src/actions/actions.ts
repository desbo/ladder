export const SET_LOGIN_MODE = 'SET_LOGIN_MODE';
export const SIGN_IN = 'SIGN_IN';
export const SIGN_OUT = 'SIGN_OUT';

export function signIn(username: string) {
  return {
    type: SIGN_IN,
    username: username
  }
}

export function setLoginMode(mode: LoginMode) {
  return {
    type: SET_LOGIN_MODE,
    mode: mode
  }
}