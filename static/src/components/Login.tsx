import * as React from 'react';

import InputField from 'components/form/InputField';
import { ReactInstance, Ref } from 'react';

export default class Login extends React.Component {
  nameInput: null | HTMLInputElement
  emailInput: null | HTMLInputElement
  passwordInput: null | HTMLInputElement

  props: { 
    input: LoginFormInput,
    mode: LoginMode,
    selectLogin: () => void,
    selectRegister: () => void,
    register: (username: string, email: string, password: string) => Promise<any>,
    signIn: (email: string, password: string) => Promise<any>
  }

  validate(f: Function) {
    const fieldsToCheck = this.props.mode === 'register' ? 
      [this.nameInput, this.emailInput, this.passwordInput] :
      [this.emailInput, this.passwordInput];

    if (fieldsToCheck.every(e => e.validity.valid)) 
      return f();
  }

  render() {
    const userInput = this.props.input;
    
    return (
      <div className="columns is-centered">
        <div className="column is-centered is-one-third">
          <div className="tabs is-centered is-boxed">
            <ul>
              <li className={this.props.mode === 'login' ? 'is-active' : ''}>
                <a onClick={this.props.selectLogin}>login</a>
              </li>

              <li className={this.props.mode === 'register' ? 'is-active' : ''}>
                <a onClick={this.props.selectRegister}>register</a>
              </li>
            </ul>
          </div>


          <form onSubmit={e => e.preventDefault()}>
            {this.props.mode === 'register' && 
              <InputField label="player name" type="text" required 
                inputRef={(e: HTMLInputElement) => this.nameInput = e} />
            }

            <InputField label="email" type="email" required 
              inputRef={(e: HTMLInputElement) => this.emailInput = e} />

            <InputField label="password" type="password" required 
              inputRef={(e: HTMLInputElement) => this.passwordInput = e} />

            <div className="buttons is-centered">
              {this.props.mode === 'login' ?
                <button 
                  type="submit"
                  className="button is-primary is-medium" 
                  onClick={() => this.validate(() => 
                    this.props.signIn(userInput.email, userInput.password))
                  }>login</button> :

                <button className="button is-primary is-medium" 
                  onClick={() => this.validate(() => {
                    this.props.register(userInput.username, userInput.email, userInput.password);
                  })}>register</button>
              }
            </div>
          </form>
        </div>
      </div>
    )
  }
}