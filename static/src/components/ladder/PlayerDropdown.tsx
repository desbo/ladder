import * as React from 'react';

type Props = {
  players: Array<LadderPlayer>,
  selected: LadderPlayer,
  onSelect: (p: LadderPlayer) => any
}

class PlayerDropdown extends React.Component<Props> {
  constructor(props: Props) {
    super(props)
  }

  select(event: React.ChangeEvent<HTMLSelectElement>) {
    const player = this.props.players.find(p => p.name === event.target.value);

    if (player) {
      this.props.onSelect(player);
    }
  }

  render() {
    const selectValue = this.props.selected ? this.props.selected.name : "empty";

    return (
      <div className="select is-pulled-right">
        <select value={selectValue} onChange={e => this.select(e)}>
          <option value="empty" className="is-hidden">opponent</option>
          {this.props.players.sort((a, b) => a.name.localeCompare(b.name)).map(p => 
            <option value={p.name} key={p.name}>{p.name}</option>
          )}
        </select>
      </div>
    )
  }
}

export default PlayerDropdown;