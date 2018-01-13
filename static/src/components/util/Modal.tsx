import * as React from 'react';

type Props = {active: boolean, level: ModalMessageLevel, close: Function}

export default class Modal extends React.Component<Props> {
  constructor(props: Props) {
    super(null)
  }

  render() {
    return (
      <div className={`modal ${this.props.active ? 'is-active' : ''}`}>
        <div className="modal-background"></div>
        <div className="modal-content">
          {this.props.children}
        </div>
        <button onClick={() => this.props.close()} 
                className="modal-close is-large" 
                aria-label="close"></button>
      </div>
    )
  }
}