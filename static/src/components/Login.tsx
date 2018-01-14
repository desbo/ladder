import * as React from 'react';

import { connect, Dispatch } from 'react-redux';
import { ReactInstance, Ref, ChangeEvent } from 'react';

import { Actions , setLoginMode, showErrorModal, showInfoModal } from 'actions/actions';
import { firebase } from 'auth';
import { User } from 'firebase';

import InputField from 'components/form/InputField';

const mapStateToProps = (state: AppState) => ({
  user: state.user,
  mode: state.view.loginMode
});

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  setLoginMode: (mode: LoginMode) => 
    dispatch(setLoginMode(mode)),

  signIn: (email: string, password: string): Promise<any> => 
    firebase.signIn(email, password),

  register: (username: string, email: string, password: string): Promise<any> =>
    firebase.register(username, email, password)
      .then((user: User) => dispatch({
        type: Actions.SIGN_IN,
        username: user.displayName
      })),

  userFormInput: (field: string, value: string) => dispatch({
    type: Actions.USER_FORM_INPUT,
    field,
    value
  }),

  showErrorModal: (error: string) => dispatch(showErrorModal(error))
});

type Props = { 
  user: UserState,
  mode: LoginMode,
  setLoginMode: (mode: LoginMode) => any,
  signIn: (email: string, password: string) => Promise<any>,
  register: (username: string, email: string, password: string) => Promise<any>,
  userFormInput: (field: string, value: string) => any,
  showErrorModal: (error: string) => any,
  fullWidth?: boolean,
  registerOnly?: boolean,
  onRegister?: () => any,
}

class Login extends React.Component<Props> {
  nameInput: null | HTMLInputElement
  emailInput: null | HTMLInputElement
  passwordInput: null | HTMLInputElement

  public static defaultProps: Partial<Props> = {
    fullWidth: false,
    registerOnly: false,
    onRegister: () => null
  }

  componentDidMount() {
    if (this.props.registerOnly) {
      this.props.setLoginMode('register');
    }
  }

  validate(f: Function) {
    const fieldsToCheck = this.props.mode === 'register' ? 
      [this.nameInput, this.emailInput, this.passwordInput] :
      [this.emailInput, this.passwordInput];

    if (fieldsToCheck.every(e => e.validity.valid)) 
      return f();
  }

  onError(e: AppError) {
    return this.props.showErrorModal(e.message);
  }

  inputEmail(email: string) {
    return this.props.userFormInput('email', email);
  }

  inputUsername(username: string) {
    return this.props.userFormInput('username', username);
  }

  inputPassword(password: string) {
    return this.props.userFormInput('password', password);
  }

  render() {
    const userInput = this.props.user.formInput;

    return (
      <div className="columns is-centered">
        <div className={`column is-centered ${this.props.fullWidth ? '' : 'is-one-third'}`}>
          {!this.props.registerOnly && 
            <div className="tabs is-centered is-boxed">
              <ul>
                <li className={this.props.mode === 'login' ? 'is-active' : ''}>
                  <a onClick={() => this.props.setLoginMode('login')}>login</a>
                </li>

                <li className={this.props.mode === 'register' ? 'is-active' : ''}>
                  <a onClick={() => this.props.setLoginMode('register')}>register</a>
                </li>
              </ul>
            </div>
          }

          <form onSubmit={e => e.preventDefault()}>
            {this.props.mode === 'register' && 
              <InputField label="name" type="text" required 
                value={userInput.username || ''}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.inputUsername(e.target.value)}
                inputRef={(e: HTMLInputElement) => this.nameInput = e} />
            }

            <InputField label="email" type="email" required 
              value={userInput.email || ''}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.inputEmail(e.target.value)}
              inputRef={(e: HTMLInputElement) => this.emailInput = e} />

            <InputField label="password" type="password" required
              value={userInput.password || ''}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.inputPassword(e.target.value)}
              inputRef={(e: HTMLInputElement) => this.passwordInput = e} />

            <div className="buttons is-centered">
              {this.props.mode === 'login' ?
                <button 
                  type="submit"
                  className="button is-primary is-medium" 
                  onClick={() => this.validate(() => 
                    this.props.signIn(userInput.email, userInput.password))
                      .catch(this.onError.bind(this))
                  }>login</button> :

                <button className="button is-primary is-medium" 
                  onClick={() => this.validate(() => {
                    this.props.register(userInput.username, userInput.email, userInput.password)
                      .then(() => this.props.onRegister())
                      .catch(this.onError.bind(this))
                  })}>register</button>
              }
            </div>
          </form>
        </div>
      </div>
    )
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Login)