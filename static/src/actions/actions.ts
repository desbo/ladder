import { error } from "util";
import { Action, AnyAction } from "redux";

export enum Actions {
  SET_LOGIN_MODE = 'SET_LOGIN_MODE',
  SIGN_IN = 'SIGN_IN',
  SIGN_OUT = 'SIGN_OUT',
  USER_FORM_INPUT = 'USER_FORM_INPUT',

  SET_PLAYER_LADDERS = 'SET_PLAYER_LADDERS',

  SHOW_MODAL = 'SHOW_MODAL',
  HIDE_MODAL = 'HIDE_MODAL'
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

export function setPlayerLadders(ladders: LadderState) {
  return {
    type: Actions.SET_PLAYER_LADDERS,
    ladders: ladders
  }
}

function showModal(message: string, level: ModalMessageLevel): AnyAction {
  return {
    type: Actions.SHOW_MODAL,
    message: message,
    level: level
  }
}

export const showErrorModal = (errorMessage: string): AnyAction => 
  showModal(errorMessage, 'error')

export const showInfoModal = (message: string): AnyAction => 
  showModal(message, 'info')