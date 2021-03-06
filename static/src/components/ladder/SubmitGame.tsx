import * as React from 'react';

import API from 'api';

type Props = { 
  ladder: Ladder, 
  user: UserState, 
  opponent: LadderPlayer,
  onSubmit: (game: Game) => any
}

type State = {
  submitting: boolean,
  scores: {
    user: string,
    opponent: string
  }
}

type ScoreInputProps = { 
  name: string, 
  placeholder: string, 
  value: string, 
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void, 
  [propName: string]: any
};

const PlayerScoreInput = ({ name, placeholder, value, onChange, ...props}: ScoreInputProps) => 
  <div className="column is-5 level-item">
    <h2 className="subtitle is-5 has-text-centered">
      {name}
    </h2>
    <input 
      value={value} 
      onChange={onChange} 
      className="input has-text-centered" 
      type="number" 
      placeholder={placeholder} 
      {...props} 
      required />
  </div>

export default class SubmitGame extends React.Component<Props, State> {
  private static defaultProps: Partial<Props> = {
    onSubmit: (game: Game) => null
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      submitting: false,
      scores: {
        user: "",
        opponent: "" 
      }
    }
  }

  setScore(player: 'user' | 'opponent', score: string) {
    const newScores = Object.assign({}, this.state.scores, {
      [player]: score
    });

    this.setState({
      scores: newScores
    });
  }

  submit() {
    this.setState({ submitting: true });

    return API.submitGame(
      this.props.ladder.id,
      this.props.opponent, 
      parseInt(this.state.scores.user), 
      parseInt(this.state.scores.opponent)
    ).then((game: Game) => {
      this.setState({
        submitting: false,
        scores: {
          user: "",
          opponent: "" 
        }
      });

      this.props.onSubmit(game);
    })
  }

  render() {
    return(
      <form className="box" onSubmit={e => e.preventDefault()}>
        <div className="columns level is-mobile">
          <PlayerScoreInput 
            value={this.state.scores.user || ""}
            name={this.props.user.username}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.setScore('user', e.target.value)}
            placeholder="your score" />

          <div className="column level-item is-2">
            <h2 className="title is-5 has-text-centered">
              vs
            </h2>
          </div>
          
          <PlayerScoreInput 
            value={this.state.scores.opponent}
            name={this.props.opponent && this.props.opponent.name || "?"} 
            placeholder="their score"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.setScore('opponent', e.target.value)}
            disabled={!this.props.opponent}
            />
        </div>

        <div className="has-text-centered">
          <button disabled={this.state.submitting || !this.state.scores.user || !this.state.scores.opponent}
                  className={`button is-primary ${this.state.submitting ? 'is-loading' : ''}`}
                  onClick={this.submit.bind(this)}>submit</button>
        </div>
      </form>
    )
  }
}