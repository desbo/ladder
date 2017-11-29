import * as React from 'react';

import InputField from 'components/form/InputField';

export default class NewLadderForm extends React.Component<{className: string, createLadder: Function}, {ladderName: string}> {
  createLadder() {
    return this.props.createLadder(this.state.ladderName)
      .then(() => this.setState({
        ladderName: ''
      }))
  }

  render() {
    return <div className={this.props.className}>
      <h1 className="title is-4">create new ladder</h1>

      <form onSubmit={e => e.preventDefault()}>
        <InputField label="name" type="text" onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.setState({
          ladderName: e.target.value
        })} required />

        <div className="buttons is-centered">
          <button 
            type="submit"
            className="button is-primary" 
            onClick={() => this.createLadder()}>
              create
          </button>
        </div>
      </form>
    </div>
  }
}