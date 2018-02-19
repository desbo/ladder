import * as React from 'react';

type Props = {
  players: Array<LadderPlayer>,
  onSelect: (p: LadderPlayer) => any
}

type State = {
  active: boolean
}

class PlayerDropdown extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      active: false
    }
  }

  componentDidUpdate() {
    if (this.state.active) {
      document.addEventListener('click', () => {
        this.setState({ active: false });
      }, { once: true });
    }
  }

  select(player: LadderPlayer) {
    this.props.onSelect(player);
    this.setState({
      active: false
    })
  }

  render() {
    return (
      <div className={`dropdown is-pulled-right ${this.state.active ? 'is-active' : ''}`}>
        <div className="dropdown-trigger">
          <button className="button" aria-haspopup="true" aria-controls="dropdown-menu"
            onClick={() => this.setState({
              active: !this.state.active
            })}>

            <span>opponent</span>
            <span className="icon is-small">
              <i className="fa fa-angle-down" aria-hidden="true"></i>
            </span>
          </button>
        </div>

        <div className="dropdown-menu" id="dropdown-menu" role="menu">
          <div className="dropdown-content">
            {this.props.players.sort((a, b) => a.name.localeCompare(b.name)).map(p => 
              <a onClick={() => this.select(p)} key={p.name} className="dropdown-item">
                {p.name}
              </a>
            )}
          </div>
        </div>
      </div>
    )
  }
}

export default PlayerDropdown;