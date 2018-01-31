import * as React from 'react';
import { connect } from 'react-redux';
import { Dispatch } from 'redux';

import { Actions, showErrorModal } from 'actions/actions';
import { firebase } from 'auth';

import InlineInputField  from 'components/form/InlineInputField';

const mapDispatchToProps = (dispatch: Dispatch<any>) => ({
  userFormInput: (field: string, value: string) => dispatch({
    type: Actions.USER_FORM_INPUT,
    inline: true,
    field,
    value
  }),

  signIn: (email: string, password: string): Promise<any> => 
    firebase.signIn(email, password),

  showErrorModal: (error: string) => dispatch(showErrorModal(error))
})

const mapStateToProps = (state: AppState) => ({
  user: state.user
});

type Props = {
  user: UserState,
  userFormInput: (field: string, value: string) => any,
  showErrorModal: (error: string) => any,
  signIn: (email: string, password: string) => Promise<any>,
}

class NavLogin extends React.Component<Props> {
  emailInput: null | HTMLInputElement
  passwordInput: null | HTMLInputElement

  inputEmail(email: string) {
    return this.props.userFormInput('email', email);
  }

  inputPassword(password: string) {
    return this.props.userFormInput('password', password);
  }

  validate(f: Function) {
    const fieldsToCheck = [this.emailInput, this.passwordInput];
      
    if (fieldsToCheck.every(e => e.validity.valid)) return f();
  }

  onError(e: AppError) {
    return this.props.showErrorModal(e.message);
  }

  render() {
    return (
      <form className="navbar-item navbar-end" onSubmit={e => e.preventDefault()}>
        <div className="field is-horizontal">
          <div className="field-label is-small has-text-grey" style={{"flex-basis": "40px"}}>
            <label>already registered?</label>
          </div>
          
          <div className="field-body">
            <InlineInputField
              placeholder="email"
              value={this.props.user.inlineFormInput.email || ''}
              type="email" required
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.inputEmail(e.target.value)}
              inputRef={(e: HTMLInputElement) => this.emailInput = e} />

            <InlineInputField
              placeholder="password"
              value={this.props.user.inlineFormInput.password || ''}
              type="password" required
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.inputPassword(e.target.value)}
              inputRef={(e: HTMLInputElement) => this.passwordInput = e} />

            <button 
              type="submit"
              className="button is-primary is-small"
              onClick={() => this.validate(() => 
                this.props.signIn(this.props.user.inlineFormInput.email, this.props.user.inlineFormInput.password)
                  .catch(this.onError.bind(this))
              )}>login</button>
          </div>
        </div>
      </form>
    )
  }
}

export default connect(
  mapStateToProps, 
  mapDispatchToProps
)(NavLogin)
