import * as React from 'react';

import InputField from 'components/form/InputField';

type Props = {className: string, createLadder: Function};

export default class NewLadderForm extends React.Component<Props, {ladderName: string}> {
  constructor(props: Props) {
    super(props);
    this.state = {
      ladderName: ''
    }
  }

  createLadder() {
    return this.props.createLadder(this.state.ladderName)
      .then(() => this.setState({ ladderName: '' }))
  }

  render() {
    return <div className={this.props.className}>
      <h2 className="title is-4">create new ladder</h2>

      <form onSubmit={e => e.preventDefault()}>
        <InputField value={this.state.ladderName || ''} label="name" type="text" onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.setState({
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