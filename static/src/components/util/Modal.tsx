import * as React from 'react';

type Props = {active: boolean, level: ModalMessageLevel, close: Function}

export default class Modal extends React.Component<Props> {
  constructor(props: Props) {
    super(props)
  }

  render() {
    return (
      <div className={`modal ${this.props.active ? 'is-active' : ''}`}>
        <div className="modal-background" onClick={() => this.props.close()}></div>
        <div className="modal-content">
          <div className={`notification ${this.props.level === 'error' ? 'is-danger' : ''}`}>
            <button onClick={() => this.props.close()} className="delete"></button>
            {this.props.children}
          </div>
        </div>
        <button onClick={() => this.props.close()} 
                className="modal-close is-large" 
                aria-label="close"></button>
      </div>
    )
  }
}