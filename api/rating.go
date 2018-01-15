package app

import (
	"context"

	"github.com/jlouis/glicko2"
	"google.golang.org/appengine/datastore"
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
	return o.Rating
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
// and returns the new rankings (winner and loser respectively)
//
// TODO: this works on a per-game basis. The glicko2 doc says it works best with a
// longer rating period of 10-15 games, so maybe this function should take a single
// Player, look up their previous ~10 games and calculate the rating based on that.
func rank(ctx context.Context, g *Game) (float64, float64, error) {
	winner, loser := g.WinnerAndLoser()

	// winnerOpponent is the loser (the opponent of the winner)
	winnerOpponent := []glicko2.Opponent{newOpponent(loser, win)}
	loserOpponent := []glicko2.Opponent{newOpponent(winner, loss)}

	wc := rankPlayer(&winner, winnerOpponent)
	lc := rankPlayer(&loser, loserOpponent)

	_, err := datastore.PutMulti(
		ctx,
		[]*datastore.Key{winner.DatastoreKey(ctx), loser.DatastoreKey(ctx)},
		[]Player{winner, loser},
	)

	if err != nil {
		return 0, 0, err
	}

	g.SetRatingChange(winner, wc)
	g.SetRatingChange(loser, lc)

	return winner.Rating, loser.Rating, nil
}

func rankPlayer(p *Player, opponents []glicko2.Opponent) float64 {
	prev := p.Rating

	p.Rating, p.RatingDeviation, p.RatingVolatility = glicko2.Rank(
		p.Rating,
		p.RatingDeviation,
		p.RatingVolatility,
		opponents,
		τ,
	)

	return p.Rating - prev
}
