package app

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type playerResult struct {
	Player       Player
	Score        int
	RatingChange float64
}

type Game struct {
	ID      string
	Date    time.Time
	Player1 playerResult
	Player2 playerResult
}

func newPlayerResult(ctx context.Context, p *Player, score int) playerResult {
	return playerResult{
		Player: *p,
		Score:  score,
	}
}

// WinnerAndLoser returns the player Keys for the winner and loser of this match, respectively
func (g *Game) WinnerAndLoser() (Player, Player) {
	if g.Player1.Score > g.Player2.Score {
		return g.Player1.Player, g.Player2.Player
	}

	return g.Player2.Player, g.Player1.Player
}

func (g *Game) SetRatingChange(p Player, change float64) error {
	if g.Player1.Player.Equals(p) {
		g.Player1.RatingChange = change
	} else if g.Player2.Player.Equals(p) {
		g.Player2.RatingChange = change
	} else {
		return fmt.Errorf("player %s not in game %s", p.FirebaseID, g.ID)
	}

	return nil
}

// NewGame creates a new game
func NewGame(ctx context.Context, p1 *Player, p2 *Player, p1score int, p2score int) *Game {
	return &Game{
		ID:      xid.New().String(),
		Date:    time.Now(),
		Player1: newPlayerResult(ctx, p1, p1score),
		Player2: newPlayerResult(ctx, p2, p2score),
	}
}
