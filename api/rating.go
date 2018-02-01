package api

import (
	"context"

	"github.com/jlouis/glicko2"
)

// http://www.glicko.net/glicko/glicko2.pdf

const τ = 0.6

type result float64

// glicko2 values for win/loss from perspective of player to be ranked
const win result = 1.0
const loss result = 0.0

// implements glicko2.Opponent
type Opponent struct {
	Player
	Result result
}

func (o Opponent) R() float64 {
	return float64(o.Rating)
}

func (o Opponent) RD() float64 {
	return o.RatingDeviation
}

func (o Opponent) Sigma() float64 {
	return o.RatingVolatility
}

func (o Opponent) SJ() float64 {
	return float64(o.Result)
}

func newOpponent(p Player, outcome result) Opponent {
	return Opponent{
		p,
		outcome,
	}
}

// Rank updates the ranking, deviation and volatility for the Players in a Game
// and returns the ranking adjustments (winner and loser respectively).
//
// TODO: this works on a per-game basis. The glicko2 doc says it works best with a
// longer rating period of 10-15 games, so maybe this function should take a single
// Player, look up their previous ~10 games and calculate the rating based on that.
func rank(ctx context.Context, winner *Player, loser *Player) (int, int) {
	// winnerOpponent is the loser (the opponent of the winner)
	winnerOpponent := []glicko2.Opponent{newOpponent(*loser, win)}
	loserOpponent := []glicko2.Opponent{newOpponent(*winner, loss)}

	wa := rankPlayer(winner, winnerOpponent)
	la := rankPlayer(loser, loserOpponent)

	return wa, la
}

// rankPlayer updates this Player's rank and returns the difference
func rankPlayer(p *Player, opponents []glicko2.Opponent) int {
	prev := p.Rating

	r, rd, rv := glicko2.Rank(
		float64(p.Rating),
		p.RatingDeviation,
		p.RatingVolatility,
		opponents,
		τ,
	)

	p.Rating, p.RatingDeviation, p.RatingVolatility = int(r), rd, rv

	return p.Rating - prev
}
