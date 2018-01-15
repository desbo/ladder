import { error } from "util";
import { Action, AnyAction } from "redux";
import Ladders from "components/Ladders";

export enum Actions {
  SET_LOGIN_MODE = 'SET_LOGIN_MODE',
  SIGN_IN = 'SIGN_IN',
  SIGN_OUT = 'SIGN_OUT',
  USER_FORM_INPUT = 'USER_FORM_INPUT',

  SET_PLAYER_LADDERS = 'SET_PLAYER_LADDERS',
  SET_CURRENT_LADDER = 'SET_CURRENT_LADDER',

  SELECT_OPPONENT = 'SELECT_OPPONENT',

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

export function setCurrentLadder(ladder: Ladder) {
  return {
    type: Actions.SET_CURRENT_LADDER,
    ladder: ladder
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

export const selectOpponent = (opponent: LadderPlayer): AnyAction => ({
  type: Actions.SELECT_OPPONENT,
  opponent: opponent
})

export const clearOpponent = (): AnyAction => ({
  type: Actions.SELECT_OPPONENT,
  opponent: undefined
})