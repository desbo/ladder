package app

import (
	"context"
	"time"

	"github.com/rs/xid"
	"google.golang.org/appengine/datastore"
)

type playerScore struct {
	PlayerKey *datastore.Key
	Score     int
}

type Game struct {
	ID      string
	Date    time.Time
	Player1 playerScore
	Player2 playerScore
}

func newPlayerScore(ctx context.Context, p *Player, score int) playerScore {
	return playerScore{
		PlayerKey: p.DatastoreKey(ctx),
		Score:     score,
	}
}

// WinnerAndLoser returns the player Keys for the winner and loser of this match, respectively
func (g *Game) WinnerAndLoser() (*datastore.Key, *datastore.Key) {
	if g.Player1.Score > g.Player2.Score {
		return g.Player1.PlayerKey, g.Player2.PlayerKey
	}

	return g.Player2.PlayerKey, g.Player1.PlayerKey
}

func NewGame(ctx context.Context, p1 *Player, p2 *Player, p1score int, p2score int) *Game {
	return &Game{
		ID:      xid.New().String(),
		Date:    time.Now(),
		Player1: newPlayerScore(ctx, p1, p1score),
		Player2: newPlayerScore(ctx, p2, p2score),
	}
}
